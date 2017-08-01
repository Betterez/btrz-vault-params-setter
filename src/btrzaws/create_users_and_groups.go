package btrzaws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

// ServiceInformation - service informaito needed to create groups and users
type ServiceInformation struct {
	ServiceName          string
	RequiredEnvironments []string
	RequiredArn          []string
	Path                 string
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

// AddServiceArn - adds an aws arn to the service request
func (si *ServiceInformation) AddServiceArn(arn string) {
	si.RequiredArn = append(si.RequiredArn, arn)
}

// GetGroupName - return the group name for the service
func (si *ServiceInformation) GetGroupName() string {
	return fmt.Sprintf("%s-Group", si.ServiceName)
}

// CreateGroupAndUsersForService - for a given service name creates a groups, users and stores the keys in the vault
func CreateGroupAndUsersForService(awsSession *session.Session, iamService *iam.IAM, serviceInfo *ServiceInformation) error {
	_, err := iamService.CreateGroup(&iam.CreateGroupInput{
		GroupName: aws.String(serviceInfo.GetGroupName()),
		Path:      &serviceInfo.Path,
	})
	if err != nil {
		return err
	}
	for _, PolicyArn := range serviceInfo.RequiredArn {
		_, err = iamService.AttachGroupPolicy(&iam.AttachGroupPolicyInput{
			GroupName: aws.String(serviceInfo.GetGroupName()),
			PolicyArn: &PolicyArn,
		})
		if err != nil {
			return err
		}
	}
	for _, environment := range serviceInfo.RequiredEnvironments {
		currentUserName := fmt.Sprintf("user-%s-%s", serviceInfo.ServiceName, environment)
		_, err = iamService.CreateUser(&iam.CreateUserInput{
			Path:     &serviceInfo.Path,
			UserName: aws.String(currentUserName),
		})
		if err != nil {
			return err
		}
		_, err = iamService.AddUserToGroup(&iam.AddUserToGroupInput{
			GroupName: aws.String(serviceInfo.GetGroupName()),
			UserName:  &currentUserName,
		})
		if err != nil {
			return err
		}
		userKeysResponse, err := iamService.CreateAccessKey(&iam.CreateAccessKeyInput{
			UserName: &currentUserName,
		})
		if err != nil {
			return err
		}
		err = addKeysToVault(environment, userKeysResponse, serviceInfo)
		if err != nil {
			return err
		}
	}

	return nil
}

func addKeysToVault(environment string, akOutput *iam.CreateAccessKeyOutput, serviceInfo *ServiceInformation) error {
	fmt.Println("adding ", *akOutput.AccessKey.AccessKeyId)
	return nil
}
