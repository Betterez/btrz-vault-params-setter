package main

import (
	"btrzaws"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

const (
	versionNumber = 1
)

func getAwsUsernames(awsSession *session.Session, iamService *iam.IAM) ([]string, error) {
	usernames := make([]string, 20)
	usersResults, err := iamService.ListUsers(&iam.ListUsersInput{})
	if err != nil {
		return nil, err
	}
	for _, userInformation := range usersResults.Users {
		usernames = append(usernames, *userInformation.UserName)
	}
	return usernames, nil
}

func main() {
	awsSession, err := btrzaws.GetAWSSession()
	if err != nil {
		fmt.Print(err, "can't get a session")
		os.Exit(1)
	}
	log.Println("session created")
	iamService := iam.New(awsSession)
	if iamService == nil {
		fmt.Println("can't create iam")
		os.Exit(1)
	}
	//keysMetaData := make([]*iam.AccessKeyMetadata, 40)
	fmt.Printf("Version %d\n", versionNumber)
	accessKeysFilter := &iam.ListAccessKeysInput{}
	accessKeysFilter.SetMaxItems(80)
	searchMask := "[\\w+=,.@-]+"
	accessKeysFilter.SetUserName(searchMask)
	policiesInput := &iam.ListUserPoliciesInput{}
	usersPoliciesNames := make([]string, 0)
	//accessKeysFilter.SetUserName("qualys2017")
	usernames, err := getAwsUsernames(awsSession, iamService)
	if err != nil {
		fmt.Print("error ", err, "exiting\n")
		os.Exit(1)
	}
	for _, username := range usernames {
		policiesInput.SetUserName(username)
		usersPolicies, _ := iamService.ListUserPolicies(policiesInput)
		for _, currentPolicy := range usersPolicies.PolicyNames {
			usersPoliciesNames = append(usersPoliciesNames, *currentPolicy)
		}
		attachedPolicies, _ := iamService.ListAttachedUserPolicies(&iam.ListAttachedUserPoliciesInput{UserName: &username})
		for _, currentPolicy := range attachedPolicies.AttachedPolicies {
			usersPoliciesNames = append(usersPoliciesNames, *currentPolicy.PolicyName)
		}
		if len(usersPoliciesNames) > 1 {
			fmt.Printf("Policy listing for %s\n%v\n\n", username, usersPoliciesNames)
			usersPoliciesNames = make([]string, 0)
		}
		// accessKeysFilter.SetUserName(*userInformation.UserName)
		// keyResults, err := iamService.ListAccessKeys(accessKeysFilter)
		// if err != nil {
		// 	fmt.Print("aws error", err)
		// 	os.Exit(1)
		// }
		// for _, keyInfo := range keyResults.AccessKeyMetadata {
		// 	keysMetaData = append(keysMetaData, keyInfo)
		// }
	}
	//fmt.Print(keysMetaData)

}
