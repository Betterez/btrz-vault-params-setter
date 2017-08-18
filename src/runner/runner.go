package main

import (
	"btrzaws"
	_ "btrzdb"
	"btrzutils"
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
	const fileName = "./secrets/log_entries.json"
	driver, err := btrzutils.CreateConnectionFromSecretsFile(fileName)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	if !driver.IsAuthenticated() {
		fmt.Println("not authenticated!")
		os.Exit(1)
	}
	fmt.Println("account name:", driver.GetAccountName())
	users, err := driver.GetUsers()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	for _, user := range users {
		fmt.Println(user)
	}
}
