package btrzaws

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

// ServiceInformation - service informaito needed to create groups and users
type ServiceInformation struct {
	ServiceName          string
	RequiredEnvironments []string
	RequiredArn          []string
}

// GenerateServiceInformation - create a ServiceInformation with default settings
func GenerateServiceInformation(serviceName string) *ServiceInformation {
	return &ServiceInformation{
		ServiceName:          serviceName,
		RequiredEnvironments: []string{"staging", "sandbox", "production"},
		RequiredArn:          []string{},
	}
}

// CreateGroupAndUsersForService - for a given service name creates a groups, users and stores the keys in the vault
func CreateGroupAndUsersForService(awsSession *session.Session, iamService *iam.IAM, serviceName string) {

}
