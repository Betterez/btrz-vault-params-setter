package main

import (
	"btrzaws"
	"btrzutils"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/service/iam"
	simplejson "github.com/bitly/go-simplejson"
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
	vaultSetings, err := btrzutils.LoadVaultInfoFromJSONFile("secrets/secrets.json", "staging")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	connection, err := btrzutils.CreateVaultConnection(vaultSetings)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dataString, _ := connection.GetJSONValue("secret/betterez-app")
	jsonData, err := simplejson.NewJson([]byte(dataString))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if _, value := jsonData.CheckGet("data"); !value {
		fmt.Print("no data for this value\n")
	} else {
		fmt.Println(jsonData.Get("data"))
	}
}
