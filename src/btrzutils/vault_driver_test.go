package btrzutils

import (
	"fmt"
	"os"
	"testing"

	"github.com/bitly/go-simplejson"
)

func GetToken() string {
	sourceFile := "../../secrets/secrets.json"
	jsonFile, err := os.Open(sourceFile)
	if err != nil {
		fmt.Print(err)
		return ""
	}
	jsonData, err := simplejson.NewFromReader(jsonFile)
	if err != nil {
		fmt.Print(err)
		return ""
	}
	fmt.Println("loading token from json")
	token, err := jsonData.Get("staging").Get("vault").Get("token").String()
	//fmt.Println(err)
	return token
}

func TestConnectionFromJSONData(t *testing.T) {
	token := GetToken()
	if token == "" {
		t.Fatal("can't get token")
	}
	driver, err := CreateVaultConnectionFromParameters("vault-staging.betterez.com", token, 9000)
	if err != nil {
		t.Fatal(err)
	}
	if driver.GetValutStatus() != "online" {
		t.Fatalf("Driver status is %s", driver.GetValutStatus())
	}
}

func TestJSONValues(t *testing.T) {
	const testPath = "secret/test_from_go"
	const testJSONData = `{"user":"jarjar binxx"}`
	token := GetToken()
	if token == "" {
		t.Fatal("can't get token")
	}
	driver, err := CreateVaultConnectionFromParameters("vault-staging.betterez.com", token, 9000)
	if err != nil {
		t.Fatal(err)
	}
	if driver.GetValutStatus() != "online" {
		t.Fatalf("Driver status is %s", driver.GetValutStatus())
	}
	code, err := driver.PutJSONValue(testPath, testJSONData)
	if err != nil {
		t.Fatal(err)
	}
	if code > 399 {
		t.Fatalf("http code error, %d was returned", code)
	}
	jsonData, err := driver.GetJSONValue(testPath)
	if err != nil {
		t.Fatal(err)
	}
	result, err := simplejson.NewJson([]byte(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	resultData, err := result.Get("data").Get("user").String()
	if err != nil {
		t.Fatal(err)
	}
	if resultData != "jarjar binxx" {
		t.Fatalf("return valuse equal to %s", resultData)
	}

}
