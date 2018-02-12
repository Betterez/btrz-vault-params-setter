package btrzdb

import (
	"fmt"
	"os"
)

// DeploymentData - deployment descriptor
type DeploymentData struct {
	ServerAddress string
	Username      string
	Password      string
	DatabaseName  string
}

// IsLegal check if the data is legal
func (data *DeploymentData) IsLegal() bool {
	if data.ServerAddress == "" {
		return false
	}
	if data.Username != "" {
		if data.Password == "" {
			return false
		}
	}
	return true
}

// IsAuthenticated - does the info contains authentication info
func (data *DeploymentData) IsAuthenticated() bool {
	if data.Username != "" && data.IsLegal() {
		return true
	}
	return false
}

// MakeDialString - create a dial string from the info
func (data *DeploymentData) MakeDialString() string {
	var connectionString string
	if data.IsAuthenticated() {
		connectionString = fmt.Sprintf("mongodb://%s:%s@%s/admin", data.Username, data.Password, data.ServerAddress)
	} else {
		connectionString = fmt.Sprintf("mongodb://%s/%s", data.ServerAddress, data.DatabaseName)
	}
	return connectionString
}

// MakeRestoreString - create a restore string
func (data *DeploymentData) MakeRestoreString() string {
	var restoreString string
	if data.IsAuthenticated() {
		restoreString = fmt.Sprintf("'mongorestore --authenticationDatabase admin -u %s -p %s'", data.Username, data.Password)
	} else {
		restoreString = "mongorestore"
	}
	return restoreString
}

// CreateDeploymentFromEnvVars pulls data from env
func CreateDeploymentFromEnvVars() *DeploymentData {
	result := &DeploymentData{
		ServerAddress: os.Getenv("MONGO_SERVER_ADDRESS"),
		Username:      os.Getenv("MONGO_USERNAME"),
		Password:      os.Getenv("MONGO_PASSWORD"),
		DatabaseName:  os.Getenv("MONGO_DATABASE_NAME"),
	}
	return result
}
