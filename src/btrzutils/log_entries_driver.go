package btrzutils

import (
	"errors"
	simplejson "github.com/bitly/go-simplejson"
	"net/http"
	"os"
	"time"
)

// LogEntriesConnection - represent a log entry api connection
type LogEntriesConnection struct {
	apiKey    string
	apiKeyID  string
	accountID string
}

const (
	// LeAPIHeader  the header string for log entries
	LeAPIHeader = "x-api-key"
)

// CreateConnection - returns new connection or an error
func CreateConnection(APIKey, APIKeyID, accountID string) (*LogEntriesConnection, error) {
	result := &LogEntriesConnection{
		apiKey:    APIKey,
		accountID: accountID,
		apiKeyID:  APIKeyID,
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

// GetUsers - list users in the account
func (conn *LogEntriesConnection) GetUsers() ([]string, error) {
	const url = "https://rest.logentries.com/management/accounts/:accountid/users"
	usersRequest, _ := http.NewRequest("GET", url, nil)
	usersRequest.Header.Set(LeAPIHeader, conn.apiKey)
	httpClient := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	response, err := httpClient.Do(usersRequest)
	if err != nil {
		return nil, err
	}
	if response.StatusCode > 399 {
		return nil, errors.New("Bad http code")
	}
	foundUsers := make([]string, 0)
	return foundUsers, nil
}
