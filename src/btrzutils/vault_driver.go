package btrzutils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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
	// AllReposPath - root path for all repo json hash
	AllReposPath = "secret/all_repos_path_storage_hash"
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
	if response.StatusCode == 404 {
		return "", nil
	}
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

// EnsureRepositoryExists - make sure that a repo by this name is listed within vault.
func (v *VaultServer) EnsureRepositoryExists(repositoryName string) (int, error) {
	if v.GetVaultStatus() != VaultOnline {
		return 0, fmt.Errorf("Vault is not online, %s", v.GetVaultStatus())
	}
	allReposString, err := v.GetJSONValue(AllReposPath)
	if err != nil {
		log.Println(err,"getting json for all repos")
		return 0, err
	}
	// not initialized yet
	if allReposString == "" {
		log.Println("new repo in this environment, creating registry")
		allReposJSON := simplejson.New()
		allReposJSON.Set("repos", map[string]int{repositoryName: 1})
		allReposJSONString, err1 := allReposJSON.MarshalJSON()
		if err1 != nil {
			return 0, err1
		}
		return v.PutJSONValue(AllReposPath, string(allReposJSONString))
	}
	//JSONToLoad := simplejson.New()
	valuesData, err := simplejson.NewJson([]byte(allReposString))
	if err != nil {
		return 0, err
	}
	reposJSONData := valuesData.Get("data")
	reposJSONData.Get("repos").Set(repositoryName, 1)
	valuesDataString, err := reposJSONData.MarshalJSON()
	if err != nil {
		return 0, err
	}
	return v.PutJSONValue(AllReposPath, string(valuesDataString))
}

// ListAllRepositories - lists all the repositories in the vault index
func (v *VaultServer) ListAllRepositories() ([]string, error) {
	if v.GetVaultStatus() != VaultOnline {
		return nil, fmt.Errorf("Vault is not online, %s", v.GetVaultStatus())
	}
	result := []string{}
	allReposString, err := v.GetJSONValue(AllReposPath)
	if err != nil {
		return nil, err
	}
	valuesData, err := simplejson.NewJson([]byte(allReposString))
	if err != nil {
		return nil, err
	}
	allReposeValues := valuesData.Get("data").Get("repos")
	if repositoriesMap, err := allReposeValues.Map(); err == nil {
		for key := range repositoriesMap {
			result = append(result, key)
		}
	}

	return result, nil
}

// RemoveRepositoryFromList - delete a repository from the list
func (v *VaultServer) RemoveRepositoryFromList(repositoryName string) (int, error) {
	if v.GetVaultStatus() != VaultOnline {
		return 0, fmt.Errorf("Vault is not online, %s", v.GetVaultStatus())
	}
	allReposString, err := v.GetJSONValue(AllReposPath)
	if err != nil {
		return 0, err
	}
	// not initialized yet
	if allReposString == "" {
		return 200, nil
	}
	reposJSONdata, err := simplejson.NewJson([]byte(allReposString))
	if err != nil {
		return 0, err
	}
	reposJSONdata = reposJSONdata.Get("data")
	if _, value := reposJSONdata.Get("repos").CheckGet(repositoryName); !value {
		return 208, nil
	}
	reposJSONdata.Get("repos").Del(repositoryName)
	if err != nil {
		return 0, err
	}
	revisedReposData, err := reposJSONdata.MarshalJSON()
	if err != nil {
		return 0, err
	}
	fmt.Println("posting", string(revisedReposData))
	return v.PutJSONValue(AllReposPath, string(revisedReposData))
}

// SetValuesForRepository - sets a value for a repository
func (v *VaultServer) SetValuesForRepository(repositoryName, values string, append bool) (int, error) {
	code, err := v.EnsureRepositoryExists(repositoryName)
	if err != nil {
		log.Print(err,"ensuring repo")
		return code, err
	}
	return v.AddValuesInPath("secret/"+repositoryName, values)
}

// GetRepositoryValues = get all values for a repository
func (v *VaultServer) GetRepositoryValues(repositoryName string) (*simplejson.Json, error) {
	if v.GetVaultStatus() != VaultOnline {
		return nil, errors.New("not connected, locked or not authenticated")
	}
	path := "secret/" + repositoryName
	existingDataString, err := v.GetJSONValue(path)
	if err != nil {
		return nil, err
	}
	existingDataJSON, err := simplejson.NewJson([]byte(existingDataString))
	if err != nil {
		return nil, err
	}
	return existingDataJSON.Get("data"), nil
}

// AddValuesInPath - adds values in the selected path without deleting other values.
// exsiting values will be overwritten
func (v *VaultServer) AddValuesInPath(path, values string) (int, error) {
	JSONToLoad := simplejson.New()
	valuesData, err := simplejson.NewJson([]byte(values))
	if err != nil {
		return 0, err
	}
	log.Println("getting current values from path", path)
	existingDataString, err := v.GetJSONValue(path)
	if err != nil {
		log.Println(err, "getting current values from path", path)
		return 0, err
	}
	if existingDataString == "" {
		log.Println("seems to be new entry, craeting empty dataset")
		existingDataString = "{}"
	}
	existingDataJSON, err := simplejson.NewJson([]byte(existingDataString))
	if err != nil {
		log.Printf("%v loading json from info \r\n", err)
		return 0, err
	}
	if vaultJSONData, exists := existingDataJSON.CheckGet("data"); exists {
		log.Println("This entry is not empty, pulling data to add new values.")
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
	}else{
		log.Println("This entry is empty, nothing to pull.")
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
	log.Println("updating new vault string...")
	code, err := v.PutJSONValue(path, string(JSONToLoadData))
	if err!=nil{
		log.Println(err,"updating new vault string!")
	}
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
