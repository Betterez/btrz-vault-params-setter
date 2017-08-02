package btrzutils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	simplejson "github.com/bitly/go-simplejson"
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
	Address string `json:"address"`
	Port    int    `json:"port"`
	Token   string `json:"token"`
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
// NOTE vault must include the path /secret/ befor.
func (v *VaultServer) GetJSONValue(path string) (string, error) {
	if v.GetVaultStatus() != VaultOnline {
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
	if v.GetVaultStatus() != VaultOnline {
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

// AddValuesInPath - adds values in the selected path without deleting other values.
func (v *VaultServer) AddValuesInPath(path, values string) (int, error) {
	JSONToLoad := simplejson.New()
	valuesData, err := simplejson.NewJson([]byte(values))
	if err != nil {
		return 0, err
	}
	existingDataString, err := v.GetJSONValue(path)
	if err != nil {
		return 0, err
	}
	existingDataJSON, err := simplejson.NewJson([]byte(existingDataString))
	if err != nil {
		return 0, err
	}
	if vaultJSONData, exists := existingDataJSON.CheckGet("data"); exists {
		allVaultKeys, err1 := vaultJSONData.Map()
		if err1 != nil {
			return 0, err1
		}
		for key := range allVaultKeys {
			keyValue, err1 := vaultJSONData.Get(key).String()
			if err1 != nil {
				return 0, err1
			}
			JSONToLoad.Set(key, keyValue)
		}
	}
	valuesDataMap, err := valuesData.Map()
	if err != nil {
		return 0, err
	}
	for key := range valuesDataMap {
		keyValue, err1 := valuesData.Get(key).String()
		if err1 != nil {
			return 0, err1
		}
		JSONToLoad.Set(key, keyValue)
	}
	JSONToLoadData, err := JSONToLoad.MarshalJSON()
	if err != nil {
		return 0, err
	}
	code, err := v.PutJSONValue(path, string(JSONToLoadData))
	return code, err
}

//GetVaultStatus - returns current vault status
func (v *VaultServer) GetVaultStatus() string {
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

// LoadVaultInfoFromJSONFile - loads vault parameters from json file
func LoadVaultInfoFromJSONFile(filename, environment string) (*VaultConnectionParameters, error) {
	jsonFile, err := os.Open(filename)
	result := &VaultConnectionParameters{}
	if err != nil {
		return nil, err
	}
	jsonData, err := simplejson.NewFromReader(jsonFile)
	if err != nil {
		return nil, err
	}
	result.Token, err = jsonData.Get(environment).Get("vault").Get("token").String()
	if err != nil {
		return nil, err
	}
	result.Address, err = jsonData.Get(environment).Get("vault").Get("address").String()
	if err != nil {
		return nil, err
	}
	result.Port, err = jsonData.Get(environment).Get("vault").Get("port").Int()
	if err != nil {
		return nil, err
	}
	return result, err
}
