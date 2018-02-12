package btrzaws

import (
	"btrzdb"
	"btrzutils"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"log"
	"strings"
)

// CreateGroupAndUsersForService - for a given service name creates a groups, users and stores the keys in the vault
func CreateGroupAndUsersForService(awsSession *session.Session, iamService *iam.IAM, serviceInfo *ServiceInformation) error {
	if serviceInfo.IsInformationOK() == false {
		log.Println("Inadequate service info")
		return errors.New("Inadequate service info")
	}
	parameterFound := false
	if serviceInfo.HasAWSInfo() {
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
	}
	for _, environment := range serviceInfo.RequiredEnvironments {
		if serviceInfo.HasAWSInfo() {
			log.Println("running user creation for environment", environment, "service", serviceInfo.ServiceName)
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
					log.Printf("aws user %s was found, skipping\n",currentUserName)
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
				log.Println("Adding aws keys to vault...")
				_, err = addAWSKeysToVault(environment, userKeysResponse, serviceInfo)
				if err != nil {
					log.Println(err,"Adding aws keys to vault!")
					return err
				}
			}
		}
		if serviceInfo.HasMongoInformation() {
			log.Println("user has mongo information create/check user...")
			created, err := createMongoUser(serviceInfo, environment)
			if err != nil {
				log.Println("error creating user, ", err)
				return err
			}
			if !created {
				log.Printf("User %s was not created, it is already exists\n", serviceInfo.GetMongoUserName())
			}
		}
		if serviceInfo.LogEntryLog {
			log.Println("log entries required, creating...")
			const fileName = "./secrets/log_entries.json"
			driver, err := btrzutils.CreateConnectionFromSecretsFile(fileName)
			if err != nil {
				return err
			}
			if !driver.IsAuthenticated() {
				return errors.New("le driver not authenticated")
			}
			serviceLog, err := driver.CreateLogIfNotPresent(serviceInfo.GetLELogName(), environment)
			if err != nil {
				if err.Error() != btrzutils.ErrorLogAlreadyExists {
					return err
				}
			}
			if !serviceLog.HasTokens() {
				return errors.New("Service has no tokens")
			}
			addServiceValuesToVault(map[string]string{"LOGENTRIES_TOKEN": serviceLog.Tokens[0]}, environment, serviceInfo)
		}
	}

	return nil
}

func createMongoUser(serviceInfo *ServiceInformation, environment string) (bool, error) {
	password := btrzutils.RandStringRunes(25)
	created, err := btrzdb.CreateUser(serviceInfo.GetMongoUserName(), password, serviceInfo.MongoSettings.DatabaseName[environment], serviceInfo.MongoSettings.DatabaseRole, environment)
	if err != nil {
		log.Println(err,"Creating mongo user.")
		return false, err
	}
	if created {
		log.Println("adding new mongo user to vault")
		code, err := addMongoKeysToVault(serviceInfo.GetMongoUserName(), password, environment, serviceInfo)
		if err != nil {
			log.Println(err, "adding new mongo user to vault.")
			return false, err
		}
		log.Printf("mongo user info pushed to vault with code %d", code)
	}
	return created, nil
}

func addAWSKeysToVault(environment string, akOutput *iam.CreateAccessKeyOutput, serviceInfo *ServiceInformation) (int, error) {
	const fileName = "secrets/secrets.json"
	params, err := btrzutils.LoadVaultInfoFromJSONFile(fileName, environment)
	if err != nil {
		return 0, err
	}
	log.Println("connecting to vault server...")
	connection, err := btrzutils.CreateVaultConnection(params)
	if err != nil {
		log.Println(err,"connecting to vault server!")
		return 0, err
	}
	awsKeysString := fmt.Sprintf(`{"aws_service_key":"%s","aws_service_secret":"%s"}`, *akOutput.AccessKey.AccessKeyId, *akOutput.AccessKey.SecretAccessKey)
	log.Println("adding value in path...")
	code, err := connection.AddValuesInPath(serviceInfo.GetVaultPath(), awsKeysString)
	if err!=nil{
		log.Println(err,"adding value in path!")
	}
	return code, err
}

func addMongoKeysToVault(username, password, environment string, serviceInfo *ServiceInformation) (int, error) {
	const fileName = "secrets/secrets.json"
	params, err := btrzutils.LoadVaultInfoFromJSONFile(fileName, environment)
	if err != nil {
		return 0, err
	}
	connection, err := btrzutils.CreateVaultConnection(params)
	if err != nil {
		return 0, err
	}
	mongoString := fmt.Sprintf(`{"mongo_db_username":"%s","mongo_db_password":"%s"}`, username, password)
	code, err := connection.SetValuesForRepository(serviceInfo.ServiceName, mongoString,true)
	//code, err := connection.AddValuesInPath(serviceInfo.GetVaultPath(), mongoString)
	return code, err
}

func createJSONStringFromKeyValues(valuePairs map[string]string) string {
	infoString := "{"
	currentValue := 0
	for key, value := range valuePairs {
		infoString += fmt.Sprintf(`"%s":"%s"`, strings.ToLower(key), value)
		currentValue++
		if currentValue < len(valuePairs) {
			infoString += ","
		}
	}
	infoString += "}"
	return infoString
}

func addServiceValuesToVault(valuePairs map[string]string, environment string, serviceInfo *ServiceInformation) (int, error) {
	const fileName = "secrets/secrets.json"
	params, err := btrzutils.LoadVaultInfoFromJSONFile(fileName, environment)
	if err != nil {
		return 0, err
	}
	connection, err := btrzutils.CreateVaultConnection(params)
	if err != nil {
		return 0, err
	}
	infoString := createJSONStringFromKeyValues(valuePairs)
	code, err := connection.AddValuesInPath(serviceInfo.GetVaultPath(), infoString)
	if err==nil{
		connection.EnsureRepositoryExists(serviceInfo.ServiceName)
	}
	return code, err
}
