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

func getUsersInformation(awsSession *session.Session, iamService *iam.IAM) ([]*btrzaws.AwsUserInfo, error) {
	usersInfo := make([]*btrzaws.AwsUserInfo, 0)
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
		currentUser := &btrzaws.AwsUserInfo{Username: username}
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

func crossKeysWithUsers(awsSession *session.Session, iamService *iam.IAM) ([]*btrzaws.AwsUserInfo, error) {
	usersInfo, err := getUsersInformation(awsSession, iamService)
	if err != nil {
		return nil, err
	}
	existingKeys, err := getUsersKeysFromCSV("/home/tal/Documents/programming/go/scanner/dump/output2017-07-26 09:15:32.341402017 -0400 EDT.csv")
	if err != nil {
		return nil, err
	}
	keysMap := make(map[string]*btrzaws.AwsUserInfo)
	for _, user := range usersInfo {
		for _, currentUserKey := range user.AccessKeys {
			keysMap[currentUserKey] = user
		}
	}
	foundUsers := make([]*btrzaws.AwsUserInfo, 0)
	for _, key := range existingKeys {
		if _, exists := keysMap[key]; exists {
			foundUsers = append(foundUsers, keysMap[key])
		}
	}
	return foundUsers, nil
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
	username := "lartoTest"
	fmt.Println("creating user")
	output, err := iamService.CreateUser(&iam.CreateUserInput{UserName: &username})
	if err != nil {
		fmt.Println(err, "exiting")
		os.Exit(1)
	}
	fmt.Println("created at ", *output.User.CreateDate)

	iamService.DeleteUser(&iam.DeleteUserInput{UserName: &username})
	//keysMetaData := make([]*iam.AccessKeyMetadata, 40)
	/*fmt.Printf("Version %d\n", versionNumber)
	usersInfo, err := crossKeysWithUsers(awsSession, iamService)
	if err != nil {
		fmt.Println(err, "terminating")
		os.Exit(1)
	}
	for _, item := range usersInfo {
		fmt.Println(item.ToString())
	}*/
}

func createTestRole() {
	policyDocument := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":["ec2.amazonaws.com"]},"Action":["sts:AssumeRole"]}]}`
	roleName := "TheTestRole2"
	policyName := "arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess"
	userName := "testUser"
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
	path := "/"
	resp, err := iamService.CreateRole(&iam.CreateRoleInput{
		AssumeRolePolicyDocument: &policyDocument,
		RoleName:                 &roleName,
		Path:                     &path,
	})
	if err != nil {
		fmt.Print(err, "exiting", "\n")
		os.Exit(1)
	}
	fmt.Println(resp.String(), "created.")
	policyResponse, err := iamService.AttachRolePolicy(&iam.AttachRolePolicyInput{
		RoleName:  &roleName,
		PolicyArn: &policyName,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(policyResponse)

	userOutput, err := iamService.CreateUser(&iam.CreateUserInput{
		UserName: &userName,
	})
	if err != nil {
		fmt.Print(err, "exiting", "\n")
		os.Exit(1)
	}
	fmt.Println(userOutput)
}

// func createTestPolicy() {
// 	policyDocument := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":["ec2.amazonaws.com"]},"Action":["sts:AssumeRole"]}]}`
// 	policyName := "vue-live-test"
// 	userName := "testUser"
// 	path := "/"
// 	policyDescription := "policy for vue test service"
// 	awsSession, err := btrzaws.GetAWSSession()
// 	if err != nil {
// 		fmt.Print(err, "can't get a session")
// 		os.Exit(1)
// 	}
// 	log.Println("session created")
// 	iamService := iam.New(awsSession)
// 	if iamService == nil {
// 		fmt.Println("can't create iam")
// 		os.Exit(1)
// 	}
// 	policyResponse, err := iamService.CreatePolicy(
// 		&iam.CreatePolicyInput{
// 			PolicyDocument: &policyDocument,
// 			Path:           &path,
// 			Description:    &policyDescription,
// 			PolicyName:     &policyName,
// 		})
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}
// 	fmt.Println(policyResponse)
// }

func createTestGroup() {
	groupName := "ABetterezTest"
	policySourceName := "arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess"
	userName := "testUser"
	path := "/"
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
	_, err = iamService.CreateGroup(&iam.CreateGroupInput{
		GroupName: &groupName,
		Path:      &path,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = iamService.AttachGroupPolicy(&iam.AttachGroupPolicyInput{
		GroupName: &groupName,
		PolicyArn: &policySourceName,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	userResponse, err := iamService.CreateUser(&iam.CreateUserInput{
		Path:     &path,
		UserName: &userName,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(*userResponse)
	_, err = iamService.AddUserToGroup(&iam.AddUserToGroupInput{
		GroupName: &groupName,
		UserName:  &userName,
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("deleting user and group now")
	iamService.DeleteUser(
		&iam.DeleteUserInput{
			UserName: &userName,
		})
	iamService.DeleteGroup(&iam.DeleteGroupInput{
		GroupName: &groupName,
	})
}
func main() {
	createTestGroup()
}

//fmt.Print(keysMetaData)
