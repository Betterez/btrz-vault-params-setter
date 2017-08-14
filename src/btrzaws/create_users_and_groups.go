package btrzaws

import (
	"btrzutils"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

// CreateGroupAndUsersForService - for a given service name creates a groups, users and stores the keys in the vault
func CreateGroupAndUsersForService(awsSession *session.Session, iamService *iam.IAM, serviceInfo *ServiceInformation) error {
	if serviceInfo.IsInformationOK() == false {
		return errors.New("Inadequate service info")
	}
	parameterFound := false
	groupsResponse, err := iamService.ListGroups(&iam.ListGroupsInput{
		PathPrefix: aws.String("/"),
	})
	if err != nil {
		return err
	}
	for _, groupInfo := range groupsResponse.Groups {
		if *groupInfo.GroupName == serviceInfo.GetGroupName() {
			parameterFound = true
			break
		}
	}
	if !parameterFound {
		_, err = iamService.CreateGroup(&iam.CreateGroupInput{
			GroupName: aws.String(serviceInfo.GetGroupName()),
			Path:      &serviceInfo.Path,
		})
		if err != nil {
			return err
		}
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
		parameterFound = false
		usersListResponse, err := iamService.ListUsers(&iam.ListUsersInput{
			PathPrefix: aws.String("/"),
		})
		if err != nil {
			return err
		}
		for _, userInformation := range usersListResponse.Users {
			if *userInformation.UserName == currentUserName {
				parameterFound = true
				break
			}
		}
		if !parameterFound {
			_, err = iamService.CreateUser(&iam.CreateUserInput{
				Path:     &serviceInfo.Path,
				UserName: aws.String(currentUserName),
			})
			if err != nil {
				return err
			}
		}
		_, err = iamService.AddUserToGroup(&iam.AddUserToGroupInput{
			GroupName: aws.String(serviceInfo.GetGroupName()),
			UserName:  &currentUserName,
		})
		if err != nil {
			return err
		}
		if !parameterFound {
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
