package btrzutils

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
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
	accountName   string
	authenticated bool
}

const (
	// LeAPIHeader  the header string for log entries
	LeAPIHeader = "x-api-key"
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

// CreateConnection - returns new connection or an error
func CreateConnection(APIKey, APIKeyID, accountID string) (*LogEntriesConnection, error) {
	result := &LogEntriesConnection{
		apiKey:    APIKey,
		accountID: accountID,
		apiKeyID:  APIKeyID,
	}
	httpClient := &http.Client{}
	httpClient.Timeout = time.Duration(time.Second * 6)
	uriString := fmt.Sprintf("management/accounts/%s", accountID)
	urlStr := fmt.Sprintf("https://rest.logentries.com/%s", uriString)
	dateString := time.Now().UTC().Format("Mon, _2 Jan 2006 15:04:05 GMT")
	//dateString := "Thu, 17 Aug 2017 20:16:24 GMT"
	requestMethod := "GET"
	request, err := http.NewRequest(requestMethod, urlStr, nil)
	if err != nil {
		return nil, err
	}
	contentType := "application/json"
	request.Header.Set("Content-Type", contentType)
	request.Header.Set("Date", dateString)
	request.Header.Set("authorization-api-key", fmt.Sprintf("%s:%s", APIKeyID, GenerateSignature(APIKey, "", contentType, dateString, requestMethod, uriString)))
	//fmt.Println("header", request.Header.Get("authorization-api-key"))
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
	result, err := CreateConnection(apiKey, apiKeyID, accountResourceID)
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
func (con *LogEntriesConnection) GetUsers() ([]string, error) {
	url := fmt.Sprintf("https://rest.logentries.com/management/accounts/%s/users", con.accountID)
	usersRequest, _ := http.NewRequest("GET", url, nil)
	usersRequest.Header.Set(LeAPIHeader, con.apiKey)
	httpClient := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	response, err := httpClient.Do(usersRequest)
	if err != nil {
		return nil, err
	}
	if response.StatusCode > 399 {
		return nil, fmt.Errorf("Bad http code - %d", response.StatusCode)
	}
	foundUsers := make([]string, 0)
	return foundUsers, nil
}
