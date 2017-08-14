package main

import (
	"btrzaws"
	_ "btrzdb"
	_ "btrzutils"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func updateGroupsAndUsers() {
	serviceFile := "./services/services.json"
	if _, err := os.Stat(serviceFile); os.IsNotExist(err) {
		fmt.Printf("file %s does not exist", serviceFile)
		os.Exit(1)
	}
	servicesBytesData, err := ioutil.ReadFile(serviceFile)
	if err != nil {
		os.Exit(1)
	}
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
	servicesData := make(map[string]btrzaws.ServiceInformation)
	err = json.Unmarshal(servicesBytesData, &servicesData)
	if err != nil {
		os.Exit(1)
	}

	for serviceKey := range servicesData {
		currentServiceInfo := servicesData[serviceKey]
		btrzaws.CreateGroupAndUsersForService(awsSession, iamService, &currentServiceInfo)
	}
}
func main() {
	updateGroupsAndUsers()
}
