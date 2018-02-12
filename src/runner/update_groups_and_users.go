package main

import (
	"btrzaws"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/service/iam"
)

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
	servicesData := make(map[string]btrzaws.ServiceInformation)
	err = json.Unmarshal(servicesBytesData, &servicesData)
	if err != nil {
		log.Printf("json error:\r\n%v \r\n",err)
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
	log.Println("Iam service created.")
	for serviceKey := range servicesData {
		currentServiceInfo := servicesData[serviceKey]
		log.Printf("Applying service %s, use mongo=%v,environments=%v",currentServiceInfo.ServiceName,currentServiceInfo.HasMongoInformation(),currentServiceInfo.RequiredEnvironments)
		btrzaws.CreateGroupAndUsersForService(awsSession, iamService, &currentServiceInfo)
	}
}
