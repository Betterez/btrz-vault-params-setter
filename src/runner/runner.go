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
	info := btrzaws.GenerateServiceInformation("test-service-1")
	err = btrzaws.CreateGroupAndUsersForService(awsSession, iamService, info)
	if err != nil {
		fmt.Println(err, "creating users and group")
	}
}
func main() {

}

//fmt.Print(keysMetaData)
