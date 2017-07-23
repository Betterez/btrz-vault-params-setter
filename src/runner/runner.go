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
	fmt.Printf("Version %d", versionNumber)
	accessKeysFilter := &iam.ListAccessKeysInput{}
	usersFilter := &iam.ListUsersInput{}
	usersFilter.SetMaxItems(80)
	accessKeysFilter.SetMaxItems(80)
	//accessKeysFilter.SetUserName("[\\w+=,.@-]+")
	accessKeysFilter.SetUserName("qualys2017")
	// usersResults, err := iamService.ListUsers(usersFilter)
	// for _, userInformation := range usersResults.Users {
	// 	fmt.Print(userInformation.String())
	// }
	keyResults, err := iamService.ListAccessKeys(accessKeysFilter)
	if err != nil {
		fmt.Print("aws error", err)
		os.Exit(1)
	}
	for _, keyInfo := range keyResults.AccessKeyMetadata {
		fmt.Print(keyInfo.GoString())
	}
}
