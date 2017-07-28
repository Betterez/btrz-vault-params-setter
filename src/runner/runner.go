package main

import (
	"btrzaws"
	"encoding/csv"
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

func processUsers(awsSession *session.Session, iamService *iam.IAM) error {
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
		return err
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
	}
	return nil

}

func setup() {
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
	processUsers(awsSession, iamService)
}

func pullUsersKeysFromCSV(filename string) ([]string, error) {
	var err error
	var record []string
	keysMap := make(map[string]int, 0)
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("The file %s does not exist or is it not accesible by the current user!", filename)
	}
	usersKeys := make([]string, 0)
	fileReader, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(fileReader)
	for recordCount := 1; err == nil; recordCount++ {
		record, err = reader.Read()
		if len(record) < 2 {
			break
		}
		if recordCount > 1 && record[1] != "" {
			keysMap[record[1]] = recordCount
		}

	}
	for value := range keysMap {
		usersKeys = append(usersKeys, value)

	}
	return usersKeys, nil
}
func main() {
	keys, err := pullUsersKeysFromCSV("/home/tal/Documents/programming/go/scanner/dump/output2017-07-26 09:15:32.341402017 -0400 EDT.csv")
	if err != nil {
		fmt.Printf("Error %v. exiting\n", err)
		os.Exit(1)
	}
	for index, key := range keys {
		fmt.Printf("%d.\t%s\n", index, key)
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
