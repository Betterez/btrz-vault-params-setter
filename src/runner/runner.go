package main

import (
	"btrzaws"
	"fmt"
	"github.com/aws/aws-sdk-go/service/iam"
	"log"
	"os"
)

const (
	versionNumber = 1
)

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
	keysMetaData := make([]*iam.AccessKeyMetadata, 40)
	fmt.Printf("Version %d", versionNumber)
	accessKeysFilter := &iam.ListAccessKeysInput{}
	usersFilter := &iam.ListUsersInput{}
	usersFilter.SetMaxItems(80)
	accessKeysFilter.SetMaxItems(80)
	searchMask := "[\\w+=,.@-]+"
	accessKeysFilter.SetUserName(searchMask)
	policiesInput := &iam.ListUserPoliciesInput{}
	//accessKeysFilter.SetUserName("qualys2017")
	usersResults, err := iamService.ListUsers(usersFilter)
	if err != nil {
		fmt.Print("aws error", err)
		os.Exit(1)
	}
	for _, userInformation := range usersResults.Users {
		policiesInput.SetUserName(*userInformation.UserName)
		fmt.Println("Policy listing for ", *userInformation.UserName)
		usersPolicies, _ := iamService.ListUserPolicies(policiesInput)
		for _, currentPolicy := range usersPolicies.PolicyNames {
			fmt.Println(*currentPolicy)
		}
		accessKeysFilter.SetUserName(*userInformation.UserName)
		keyResults, err := iamService.ListAccessKeys(accessKeysFilter)
		if err != nil {
			fmt.Print("aws error", err)
			os.Exit(1)
		}
		for _, keyInfo := range keyResults.AccessKeyMetadata {
			keysMetaData = append(keysMetaData, keyInfo)
		}
	}
	fmt.Print(keysMetaData)

}
