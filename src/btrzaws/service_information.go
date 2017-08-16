package btrzaws

import (
	"fmt"
)

// MongoInformation - service mongo information
type MongoInformation struct {
	DatabaseRole string            `json:"role"`
	DatabaseName map[string]string `json:"database_name"`
}

// APIService - general api service contains key and secret, which we
type APIService struct {
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
}

// ServiceInformation - service informaito needed to create groups and users
type ServiceInformation struct {
	ServiceName          string           `json:"service_name"`
	RequiredEnvironments []string         `json:"environments"`
	RequiredArn          []string         `json:"arns"`
	Path                 string           `json:"path"`
	MongoSettings        MongoInformation `json:"mongodb"`
	APIServices          []APIService     `json:"api_services"`
}

// GenerateServiceInformation - create a ServiceInformation with default settings
func GenerateServiceInformation(serviceName string) *ServiceInformation {
	return &ServiceInformation{
		ServiceName:          serviceName,
		RequiredEnvironments: []string{"staging", "sandbox", "production"},
		RequiredArn:          []string{},
		Path:                 "/",
	}
}

// GetMongoUserName - returns mongo username for this service
func (si *ServiceInformation) GetMongoUserName() string {
	return "mongo_" + si.ServiceName
}

// HasAPIKeys - is this service info has other api keys
func (si *ServiceInformation) HasAPIKeys() bool {
	return len(si.APIServices) > 0
}

// AddServiceArn - adds an aws arn to the service request
func (si *ServiceInformation) AddServiceArn(arn string) {
	si.RequiredArn = append(si.RequiredArn, arn)
}

// GetVaultPath - return the vault (+secret) path for this service
func (si *ServiceInformation) GetVaultPath() string {
	return "secret/" + si.ServiceName
}

// GetGroupName - return the group name for the service
func (si *ServiceInformation) GetGroupName() string {
	return fmt.Sprintf("%s-Group", si.ServiceName)
}

// HasMongoInformation - returns true if this service contain mongo info
func (si *ServiceInformation) HasMongoInformation() bool {
	if si.MongoSettings.DatabaseRole != "" {
		if len(si.MongoSettings.DatabaseName) > 0 {
			return true
		}

	}
	return false
}

// IsInformationOK - Checks if the informaito provided is OK to process
func (si *ServiceInformation) IsInformationOK() bool {
	if si.ServiceName == "" {
		return false
	}
	if si.RequiredArn == nil || si.RequiredEnvironments == nil {
		return false
	}

	if len(si.RequiredEnvironments) == 0 {
		return false
	}
	if len(si.RequiredArn) == 0 {
		return false
	}
	return true
}
