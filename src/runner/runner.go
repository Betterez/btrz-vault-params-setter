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
	"time"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/bsphere/le_go"
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
	// user, err := driver.CreateUser("zz"+btrzutils.RandStringRunes(5), "zz"+btrzutils.RandStringRunes(5), "zz"+btrzutils.RandStringRunes(5)+"@betterez.com")
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// fmt.Println("user was created", user)
	// logsData, err := driver.ListLogsSet()
	// if err != nil {
	// 	fmt.Print(err)
	// 	os.Exit(1)
	// }
	// for _, logData := range logsData.Logsets {
	// 	fmt.Println(logData)
	// }
	log, err := driver.CreateNewLog("testlog", "staging")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	if !log.HasTokens() {
		fmt.Println("Log has no tokens")
	}
	if log.HasTokens() {
		fmt.Printf("log has tokens, %s\n", log.Tokens[0])
		le, err := le_go.Connect(log.Tokens[0])
		if err != nil {
			fmt.Printf("err %v while posting to le\n", err)
		}
		fmt.Println("connected to LE, sleeping 10 seconds to post...")
		counter := 0
		for {
			time.Sleep(10 * time.Second)
			le.Println(counter, "new log ", time.Now().UTC().Format("Mon, _2 Jan 2006 15:04:05 GMT"))
			fmt.Print("posted.\n")
			counter++
		}

	}
}
