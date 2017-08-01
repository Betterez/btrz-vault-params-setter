package btrzutils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

//VaultServer - a vault driver implementation
type VaultServer struct {
	address     string
	port        int
	locked      bool
	authorized  bool
	token       string
	online      bool
	initialized bool
}

const (
	// VaultOnline - vault is online
	VaultOnline = "online"
)

// VaultConnectionParameters - VaultConnectionParameters values needed to create a vault connection
type VaultConnectionParameters struct {
	Address string
	Port    int
	Token   string
}

//GetServerAddress - returns server address
func (v *VaultServer) GetServerAddress() string {
	return v.address
}

// GetServerPort - returns server port
func (v *VaultServer) GetServerPort() int {
	return v.port
}

// IsLocked - is vault locked
func IsLocked(v *VaultServer) bool {
	return v.locked
}

// CreateVaultConnection - create connection from the connection struct
func CreateVaultConnection(params *VaultConnectionParameters) (*VaultServer, error) {
	return CreateVaultConnectionFromParameters(params.Address, params.Token, params.Port)
}

// CreateVaultConnectionFromParameters - return new connection
func CreateVaultConnectionFromParameters(address, token string, port int) (*VaultServer, error) {
	result := &VaultServer{
		address: address,
		port:    port,
		token:   token,
	}
	httpClient := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	requestURL := fmt.Sprintf("http://%s:%d/v1/%s", result.address, result.port, RandStringRunes(50))
	req, _ := http.NewRequest("GET", requestURL, nil)
	req.Header.Set("X-Vault-Token", result.token)

	response, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if response.StatusCode == 503 {
		result.locked = true
	} else if response.StatusCode == 502 || response.StatusCode == 504 {
		result.online = false
	} else if response.StatusCode == 401 || response.StatusCode == 403 {
		result.authorized = false
	} else {
		result.authorized = true
		result.locked = false
		result.online = true
	}
	result.initialized = true
	return result, nil
}

// GetJSONValue - get json string from the server
func (v *VaultServer) GetJSONValue(path string) (string, error) {
	if v.GetValutStatus() != VaultOnline {
		return "", errors.New("not connected, locked or not authenticated")
	}
	requestURL := fmt.Sprintf("http://%s:%d/v1/%s", v.address, v.port, path)
	httpClient := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	req, _ := http.NewRequest("GET", requestURL, nil)
	req.Header.Set("X-Vault-Token", v.token)
	response, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// PutJSONValue - put json value into the specific path.
// NOTE vault must include the path /secret/ befor.
func (v *VaultServer) PutJSONValue(path, value string) (int, error) {
	if v.GetValutStatus() != VaultOnline {
		return 0, errors.New("not connected, locked or not authenticated")
	}
	requestURL := fmt.Sprintf("http://%s:%d/v1/%s", v.address, v.port, path)
	httpClient := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	request, _ := http.NewRequest("POST", requestURL, strings.NewReader(value))
	request.Header.Set("X-Vault-Token", v.token)
	response, err := httpClient.Do(request)
	if err != nil {
		return 0, err
	}

	return response.StatusCode, nil
}

//GetValutStatus - returns current vault status
func (v *VaultServer) GetValutStatus() string {
	if !v.initialized {
		return "not initizlied"
	}
	if v.locked {
		return "locked"
	}
	if !v.authorized {
		return "not authorized"
	}
	if v.online {
		return VaultOnline
	}
	return "unknown"
}
