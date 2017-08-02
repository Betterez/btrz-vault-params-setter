package btrzaws

import (
	"btrzutils"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

// ServiceInformation - service informaito needed to create groups and users
type ServiceInformation struct {
	ServiceName          string   `json:"service_name"`
	RequiredEnvironments []string `json:"environments"`
	RequiredArn          []string `json:"arns"`
	Path                 string   `json:"path"`
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

// GetVaultPath - return the vault (+secret) path for this service
func (si *ServiceInformation) GetVaultPath() string {
	return "secret/" + si.ServiceName
}

// GetGroupName - return the group name for the service
func (si *ServiceInformation) GetGroupName() string {
	return fmt.Sprintf("%s-Group", si.ServiceName)
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

// CreateGroupAndUsersForService - for a given service name creates a groups, users and stores the keys in the vault
func CreateGroupAndUsersForService(awsSession *session.Session, iamService *iam.IAM, serviceInfo *ServiceInformation) error {
	if serviceInfo.IsInformationOK() == false {
		return errors.New("Inadequate service info")
	}
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
		_, err = addKeysToVault(environment, userKeysResponse, serviceInfo)
		if err != nil {
			return err
		}
	}

	return nil
}

func addKeysToVault(environment string, akOutput *iam.CreateAccessKeyOutput, serviceInfo *ServiceInformation) (int, error) {
	const fileName = "secrets/secrets.json"
	params, err := btrzutils.LoadVaultInfoFromJSONFile(fileName, environment)
	if err != nil {
		return 0, err
	}
	connection, err := btrzutils.CreateVaultConnection(params)
	if err != nil {
		return 0, err
	}
	awsKeysString := fmt.Sprintf(`{"AWS_SERVICE_KEY":"%s","AWS_SERVICE_SECRET":"%s"}`, *akOutput.AccessKey.AccessKeyId, *akOutput.AccessKey.SecretAccessKey)
	code, err := connection.AddValuesInPath(serviceInfo.GetVaultPath(), awsKeysString)
	return code, err
}
