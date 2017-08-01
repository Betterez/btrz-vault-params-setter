package btrzaws

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

func getAwsUsernames(awsSession *session.Session, iamService *iam.IAM) ([]string, error) {
	usernames := make([]string, 0)
	usersResults, err := iamService.ListUsers(&iam.ListUsersInput{})
	if err != nil {
		return nil, err
	}
	for _, userInformation := range usersResults.Users {
		usernames = append(usernames, *userInformation.UserName)
	}
	return usernames, nil
}

func getUsersInformation(awsSession *session.Session, iamService *iam.IAM) ([]*AwsUserInfo, error) {
	usersInfo := make([]*AwsUserInfo, 0)
	accessKeysFilter := &iam.ListAccessKeysInput{}
	accessKeysFilter.SetMaxItems(80)
	searchMask := "[\\w+=,.@-]+"
	accessKeysFilter.SetUserName(searchMask)
	policiesInput := &iam.ListUserPoliciesInput{}
	//accessKeysFilter.SetUserName("qualys2017")
	usernames, err := getAwsUsernames(awsSession, iamService)
	if err != nil {
		fmt.Print("error ", err, "exiting\n")
		return nil, err
	}

	for _, username := range usernames {
		currentUser := &AwsUserInfo{Username: username}
		accessKeys, err := iamService.ListAccessKeys(&iam.ListAccessKeysInput{UserName: &username})
		if err != nil {
			return nil, err
		}
		for _, accessKey := range accessKeys.AccessKeyMetadata {
			currentUser.AccessKeys = append(currentUser.AccessKeys, *accessKey.AccessKeyId)
		}
		policiesInput.SetUserName(username)
		usersPolicies, _ := iamService.ListUserPolicies(policiesInput)
		for _, currentPolicy := range usersPolicies.PolicyNames {
			currentUser.Policies = append(currentUser.Policies, *currentPolicy)
		}
		attachedPolicies, _ := iamService.ListAttachedUserPolicies(&iam.ListAttachedUserPoliciesInput{UserName: &username})
		for _, currentPolicy := range attachedPolicies.AttachedPolicies {
			currentUser.Policies = append(currentUser.Policies, *currentPolicy.PolicyName)
		}
		// if len(currentUser.Policies) > 1 {
		// 	fmt.Printf("Policy listing for %s\n%v\n\n", username, currentUser.Policies)
		// }
		usersInfo = append(usersInfo, currentUser)
	}
	return usersInfo, nil

}

func getUsersKeysFromCSV(filename string) ([]string, error) {
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

func crossKeysWithUsers(awsSession *session.Session, iamService *iam.IAM) ([]*AwsUserInfo, error) {
	usersInfo, err := getUsersInformation(awsSession, iamService)
	if err != nil {
		return nil, err
	}
	existingKeys, err := getUsersKeysFromCSV("/home/tal/Documents/programming/go/scanner/dump/output2017-07-26 09:15:32.341402017 -0400 EDT.csv")
	if err != nil {
		return nil, err
	}
	keysMap := make(map[string]*AwsUserInfo)
	for _, user := range usersInfo {
		for _, currentUserKey := range user.AccessKeys {
			keysMap[currentUserKey] = user
		}
	}
	foundUsers := make([]*AwsUserInfo, 0)
	for _, key := range existingKeys {
		if _, exists := keysMap[key]; exists {
			foundUsers = append(foundUsers, keysMap[key])
		}
	}
	return foundUsers, nil
}
