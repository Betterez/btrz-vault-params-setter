package main

import (
	"btrzaws"
	"btrzutils"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	simplejson "github.com/bitly/go-simplejson"
)

const (
	strAWSKeyName = "aws_service_key"
	strAWSSecret  = "aws_service_secret"
)

var (
	listedUsers = make([]*iam.User, 0)
	usersByKeys = make(map[string]*iam.User, 0)
)

func updateMissingEmailInfo(environment string) {
	awsSession, err := btrzaws.GetAWSSession()
	if err != nil {
		fmt.Print(err, "can't get a session")
		os.Exit(1)
	}
	log.Println("loading groups and users")
	loadUsersInfo(awsSession)
	log.Println("done")
	vaultInfo, err := btrzutils.LoadVaultInfoFromJSONFile("secrets/secrets.json", environment)
	if err != nil {
		fmt.Print(err, "can't get a vault info")
		os.Exit(1)
	}
	log.Printf("vault info loaded, server %s",vaultInfo.Address)
	vault, err := btrzutils.CreateVaultConnection(vaultInfo)
	if err != nil {
		fmt.Print(err, "can't get vault connection")
		os.Exit(1)
	}
	log.Println("vault connected.")
	allReposNames, err := vault.ListAllRepositories()
	if err != nil {
		fmt.Print(err, "getting repos")
		os.Exit(1)
	}
	for _, repoName := range allReposNames {
		fmt.Println(repoName)
		fmt.Println("==============================")
		repoValues, err := getAWSKeyForRepo(vault, repoName)
		if err != nil {
			fmt.Print(err, "can't get getAWSKeyForRepo")
			continue
			//os.Exit(1)
		}
		if repoValues != nil {
			key, _ := repoValues.Get(strAWSKeyName).String()
			log.Printf("Checking %s with %s for ses...",repoName,key)
			hasAnSESEntry, err := isAWSUserHasSES(awsSession, key)
			if err != nil {
				fmt.Print(err, "can't get getAWSKeyForRepo")
				os.Exit(1)
			}
			if hasAnSESEntry {
				fmt.Printf("key %s is a ses key\n", key)
				awsSecret, err := repoValues.Get(strAWSSecret).String()
				if err != nil {
					continue
				}
				smtpPassword, _ := btrzaws.GenerateSMTPPasswordFromSecret(awsSecret)
				code, err := vault.SetValuesForRepository(repoName, fmt.Sprintf(`{"email_client_username":"%s","email_client_password":"%s"}`, key, smtpPassword), true)
				if err != nil {
					log.Printf("error %v while writing smtp values to vault, for %s repository", err, repoName)
				} else {
					log.Printf("Writing completed with code %d", code)
				}
			} else {
				fmt.Printf("key %s is not a ses key\n", key)
			}
		} else {
			fmt.Printf("doesn't contain aws key\n")
		}
		fmt.Printf("\n\n")
	}
}

func loadUsersInfo(awsSession *session.Session) error {
	if len(usersByKeys) > 0 {
		return nil
	}
	iamService := iam.New(awsSession)
	if len(listedUsers) < 1 {
		//fmt.Println("loading users")
		output, err := iamService.ListUsers(&iam.ListUsersInput{
			PathPrefix: aws.String("/"),
		})
		if err != nil {
			return err
		}
		listedUsers = output.Users
	}

	for _, awsUser := range listedUsers {
		userAccessKeys, err := iamService.ListAccessKeys(&iam.ListAccessKeysInput{
			UserName: awsUser.UserName,
		})
		if err != nil {
			continue
		}
		//fmt.Println(*awsUser.UserName)
		for _, accessKey := range userAccessKeys.AccessKeyMetadata {
			usersByKeys[*accessKey.AccessKeyId] = awsUser
		}
	}
	return nil
}

func isAWSUserHasSES(awsSession *session.Session, awsKey string) (bool, error) {
	iamService := iam.New(awsSession)
	if userInfo, ok := usersByKeys[awsKey]; ok {
		groups, err := iamService.ListGroupsForUser(&iam.ListGroupsForUserInput{
			UserName: userInfo.UserName,
		})
		if err != nil {
			return false, err
		}
		for _, group := range groups.Groups {
			policies, err := iamService.ListAttachedGroupPolicies(&iam.ListAttachedGroupPoliciesInput{
				GroupName: group.GroupName,
			})
			if err != nil {
				fmt.Println(err)
				return false, err
			}
			for _, groupPolicy := range policies.AttachedPolicies {
				if strings.Contains(*groupPolicy.PolicyName, "SES") {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

func getAWSKeyForRepo(vault *btrzutils.VaultServer, repoName string) (*simplejson.Json, error) {
	repoValues, err := vault.GetRepositoryValues(repoName)
	if err != nil {
		return nil, err
	}
	if _, ok := repoValues.CheckGet(strAWSKeyName); ok {
		return repoValues, nil
	}
	if _, ok := repoValues.CheckGet(strings.ToUpper(strAWSKeyName)); ok {
		return repoValues, nil
	}
	return nil, nil
}

func loadAWSGroups() {
	awsSession, err := btrzaws.GetAWSSession()
	if err != nil {
		fmt.Print(err, "can't get a session")
		os.Exit(1)
	}
	iamService := iam.New(awsSession)
	if iamService == nil {
		fmt.Println("can't create iam")
		os.Exit(1)
	}
	groups, err := iamService.ListGroups(&iam.ListGroupsInput{
		PathPrefix: aws.String("/"),
	})
	if err != nil {
		fmt.Print(err, "can't get the groups")
		os.Exit(1)
	}
	for _, group := range groups.Groups {
		fmt.Println(group.GoString())
	}
}
