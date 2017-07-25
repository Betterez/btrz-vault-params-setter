package main

import (
	"btrzaws"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/service/iam"
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
	//keysMetaData := make([]*iam.AccessKeyMetadata, 40)
	fmt.Printf("Version %d\n", versionNumber)
	accessKeysFilter := &iam.ListAccessKeysInput{}
	usersFilter := &iam.ListUsersInput{}
	usersFilter.SetMaxItems(80)
	accessKeysFilter.SetMaxItems(80)
	searchMask := "[\\w+=,.@-]+"
	accessKeysFilter.SetUserName(searchMask)
	policiesInput := &iam.ListUserPoliciesInput{}
	usersPoliciesNames := make([]string, 0)
	//accessKeysFilter.SetUserName("qualys2017")
	usersResults, err := iamService.ListUsers(usersFilter)
	if err != nil {
		fmt.Print("aws error", err)
		os.Exit(1)
	}
	for _, userInformation := range usersResults.Users {
		policiesInput.SetUserName(*userInformation.UserName)
		usersPolicies, _ := iamService.ListUserPolicies(policiesInput)
		for _, currentPolicy := range usersPolicies.PolicyNames {
			usersPoliciesNames = append(usersPoliciesNames, *currentPolicy)
		}
		attachedPolicies, _ := iamService.ListAttachedUserPolicies(&iam.ListAttachedUserPoliciesInput{UserName: userInformation.UserName})
		for _, currentPolicy := range attachedPolicies.AttachedPolicies {
			usersPoliciesNames = append(usersPoliciesNames, *currentPolicy.PolicyName)
		}
		if len(usersPoliciesNames) > 1 {
			fmt.Printf("Policy listing for %s\n%v\n\n", *userInformation.UserName, usersPoliciesNames)
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
