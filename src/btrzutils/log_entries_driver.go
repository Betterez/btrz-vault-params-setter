package btrzutils

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"

	"encoding/json"
	"errors"
	"fmt"
	simplejson "github.com/bitly/go-simplejson"
	"net/http"
	"os"
	"time"
)

// LogEntriesConnection - represent a log entry api connection
type LogEntriesConnection struct {
	apiKey        string
	apiKeyID      string
	accountID     string
	accountRWKey  string
	accountName   string
	authenticated bool
}

const (
	// LeAPIHeader  the header string for log entries
	LeAPIHeader = "x-api-key"
	// LERestURL  Root url
	LERestURL = "https://rest.logentries.com/"
)

// CreateSha256 - creates sha 256 from a given string
func CreateSha256(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}

// GenerateSignature - generate hash signature
func GenerateSignature(apiKey, body, contentType, dateString, requestMethod, queryPath string) string {
	encodedBodyHash := CreateSha256(body)
	canonicalString := requestMethod + contentType + dateString + queryPath + encodedBodyHash
	mac := hmac.New(sha1.New, []byte(apiKey))
	mac.Write([]byte(canonicalString))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func (con *LogEntriesConnection) setRequestHeader(request *http.Request, requestMethod, uriString, requestBody string) {
	dateString := time.Now().UTC().Format("Mon, _2 Jan 2006 15:04:05 GMT")
	contentType := "application/json"
	request.Header.Set("Content-Type", contentType)
	request.Header.Set("Date", dateString)
	request.Header.Set("authorization-api-key", fmt.Sprintf("%s:%s", con.apiKeyID, GenerateSignature(con.apiKey, requestBody, contentType, dateString, requestMethod, uriString)))
}

// CreateConnection - returns new connection or an error
func CreateConnection(APIKey, APIKeyID, accountID, rwKey string) (*LogEntriesConnection, error) {
	result := &LogEntriesConnection{
		apiKey:       APIKey,
		accountID:    accountID,
		apiKeyID:     APIKeyID,
		accountRWKey: rwKey,
	}
	httpClient := &http.Client{}
	httpClient.Timeout = time.Duration(time.Second * 6)
	uriString := fmt.Sprintf("management/accounts/%s", accountID)
	urlStr := fmt.Sprintf("%s%s", LERestURL, uriString)
	requestMethod := http.MethodGet
	request, err := http.NewRequest(requestMethod, urlStr, nil)
	if err != nil {
		return nil, err
	}
	result.setRequestHeader(request, requestMethod, uriString, "")
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode < 400 {
		result.authenticated = true
		jsonData, err := simplejson.NewFromReader(response.Body)
		defer response.Body.Close()
		if err != nil {
			return result, err
		}
		result.accountName, _ = jsonData.Get("account").Get("name").String()
	}
	return result, nil
}

// CreateConnectionFromSecretsFile - create connection from a secret file.
func CreateConnectionFromSecretsFile(fileName string) (*LogEntriesConnection, error) {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return nil, err
	}
	fileHandle, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	jsonData, err := simplejson.NewFromReader(fileHandle)
	if err != nil {
		return nil, err
	}
	accountResourceID, err := jsonData.Get("account_resource_id").String()
	if err != nil {
		return nil, err
	}
	apiKeyID, err := jsonData.Get("api_key_id").String()
	if err != nil {
		return nil, err
	}
	apiKey, err := jsonData.Get("api_key").String()
	if err != nil {
		return nil, err
	}
	rwKey, err := jsonData.Get("account_key").String()
	if err != nil {
		return nil, err
	}
	result, err := CreateConnection(apiKey, apiKeyID, accountResourceID, rwKey)
	return result, err
}

// IsAuthenticated - is this account authenticated
func (con *LogEntriesConnection) IsAuthenticated() bool {
	return con.authenticated
}

// GetAccountName - returns the account name
func (con *LogEntriesConnection) GetAccountName() string {
	return con.accountName
}

// GetUsers - list users in the account
func (con *LogEntriesConnection) GetUsers() ([]LogEntryUser, error) {
	uriString := fmt.Sprintf("management/accounts/%s/users", con.accountID)
	urlStr := fmt.Sprintf("%s%s", LERestURL, uriString)
	requestMethod := http.MethodGet
	request, _ := http.NewRequest(requestMethod, urlStr, nil)

	httpClient := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	con.setRequestHeader(request, requestMethod, uriString, "")
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode > 399 {
		return nil, fmt.Errorf("Bad http code - %d", response.StatusCode)
	}
	defer response.Body.Close()
	//bodyData, _ := ioutil.ReadAll(response.Body)
	users := &usersResponse{}
	decoder := json.NewDecoder(response.Body)
	decoder.Decode(users)
	//
	return users.Users, nil
}

// CreateUser - create new user in le
func (con *LogEntriesConnection) CreateUser(firstname, lastname, email string) (*LogEntryUser, error) {
	uriString := fmt.Sprintf("management/accounts/%s/users", con.accountID)
	urlStr := fmt.Sprintf("%s%s", LERestURL, uriString)
	requestMethod := http.MethodPost
	requestBody := fmt.Sprintf(`{
        "user":{
            "email": "%s",
            "first_name": "%s",
            "last_name": "%s"
        }
    }`, email, firstname, lastname)
	request, _ := http.NewRequest(requestMethod, urlStr, bytes.NewBuffer([]byte(requestBody)))
	httpClient := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	con.setRequestHeader(request, requestMethod, uriString, requestBody)
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode > 399 {
		return nil, fmt.Errorf("User creation returned %d error", response.StatusCode)
	}
	jd := json.NewDecoder(response.Body)
	result := &UserResponse{}
	jd.Decode(result)
	return &result.User, nil
}

// CreateNewLog - create a new log in the logset
func (con *LogEntriesConnection) CreateNewLog(logName, logSetName string) (*LogsEntriesLog, error) {
	const requestBodyTemplate = `{
	  "log": {
	    "name": "%s",
	    "structures": [],
	    "user_data": {
	      "le_agent_filename": "",
	      "le_agent_follow": "false"
	    },
	    "source_type": "token",
	    "token_seed": null,
	    "logsets_info": [
	      {
	        "id": "%s"
	      }
	    ]
	  }
	}`

	logSetsResponse, err := con.ListLogsSet()
	if err != nil {
		return nil, err
	}
	logSetID := ""
	for _, logSetInfo := range logSetsResponse.Logsets {
		if logSetInfo.Name == logSetName {
			logSetID = logSetInfo.ID
			break
		}
	}
	if logSetID == "" {
		return nil, errors.New("Log set does not exist")
	}
	requestBody := fmt.Sprintf(requestBodyTemplate, logName, logSetID)
	uriString := "management/logs"
	urlStr := fmt.Sprintf("%s%s", LERestURL, uriString)
	requestMethod := http.MethodPost
	request, _ := http.NewRequest(requestMethod, urlStr, bytes.NewBuffer([]byte(requestBody)))
	httpClient := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	con.setRequestHeader(request, requestMethod, uriString, requestBody)
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode > 399 {
		return nil, fmt.Errorf("Error code %d returned from the server.", response.StatusCode)
	}
	result := &LogsEntriesLogResponse{}
	dec := json.NewDecoder(response.Body)
	dec.Decode(result)
	return &result.Log, nil
}

// ListLogs - show current logs
func (con *LogEntriesConnection) ListLogs() (*LogEntriesLogsResponse, error) {
	uriString := "management/logs"
	urlStr := fmt.Sprintf("%s%s", LERestURL, uriString)
	requestMethod := http.MethodGet
	request, _ := http.NewRequest(requestMethod, urlStr, nil)
	httpClient := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	con.setRequestHeader(request, requestMethod, uriString, "")
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode > 399 {
		return nil, fmt.Errorf("Bad http code - %d", response.StatusCode)
	}
	defer response.Body.Close()
	dec := json.NewDecoder(response.Body)
	result := &LogEntriesLogsResponse{}
	dec.Decode(result)
	return result, nil
}

// ListLogsSet - list all log sets and their info
func (con *LogEntriesConnection) ListLogsSet() (*LogEntriesLogSetResponse, error) {
	uriString := "management/logsets"
	urlStr := fmt.Sprintf("%s%s", LERestURL, uriString)
	requestMethod := http.MethodGet
	request, _ := http.NewRequest(requestMethod, urlStr, nil)
	httpClient := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	con.setRequestHeader(request, requestMethod, uriString, "")
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode > 399 {
		return nil, fmt.Errorf("Bad http code - %d", response.StatusCode)
	}
	defer response.Body.Close()
	dec := json.NewDecoder(response.Body)
	result := &LogEntriesLogSetResponse{}
	dec.Decode(result)
	return result, nil
}
