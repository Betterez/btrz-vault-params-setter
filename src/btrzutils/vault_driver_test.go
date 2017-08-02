package btrzutils

import (
	"os"
	"testing"

	"github.com/bitly/go-simplejson"
)

func GetToken() (string, error) {
	sourceFile := "../../secrets/secrets.json"
	jsonFile, err := os.Open(sourceFile)
	if err != nil {
		return "", nil
	}
	jsonData, err := simplejson.NewFromReader(jsonFile)
	if err != nil {
		return "", err
	}
	token, err := jsonData.Get("staging").Get("vault").Get("token").String()
	return token, err
}

func TestConnectionFromJSONData(t *testing.T) {
	token, err := GetToken()
	if err != nil {
		t.Fatal(err)
	}
	if token == "" {
		t.SkipNow()
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
	token, err := GetToken()
	if err != nil {
		t.Fatal(err)
	}
	if token == "" {
		t.SkipNow()
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

func TestGetRepoValue(t *testing.T) {
	const (
		applicationPath = "secret/betterez-app"
		storagePath     = "secret/go-app-test"
		AwsKeyName      = "AWS-ACCESS-KEY"
		AwsKeyValue     = "1234567891234567"
	)
	token, err := GetToken()
	if err != nil {
		t.Fatal(err)
	}
	if token == "" {
		t.SkipNow()
	}
	driver, err := CreateVaultConnectionFromParameters("vault-staging.betterez.com", token, 9000)
	if err != nil {
		t.Fatal(err)
	}
	if driver.GetValutStatus() != "online" {
		t.Fatalf("Driver status is %s", driver.GetValutStatus())
	}
	vaultData, err := driver.GetJSONValue(applicationPath)
	if err != nil {
		t.Fatal(err)
	}
	jsonData, err := simplejson.NewJson([]byte(vaultData))
	if err != nil {
		t.Fatal(err)
	}
	appJSONData := jsonData.Get("data")
	appJSONData.Set(AwsKeyName, AwsKeyValue)
	bytesJSON, err := appJSONData.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	code, err := driver.PutJSONValue(storagePath, string(bytesJSON))
	if err != nil {
		t.Fatal(err)
	}
	if code >= 400 {
		t.Fatalf("code %d returned from server", code)
	}
	vaultData, err = driver.GetJSONValue(storagePath)
	if err != nil {
		t.Fatal(err)
	}
	jsonData, err = simplejson.NewJson([]byte(vaultData))
	if err != nil {
		t.Fatal(err)
	}
	appJSONData = jsonData.Get("data")
	extract, err := appJSONData.Get(AwsKeyName).String()
	if err != nil {
		t.Fatal(appJSONData.Get(AwsKeyName), err)
	}
	if extract != AwsKeyValue {
		t.Fatalf("extracted key equals to %s, not to %s", extract, AwsKeyValue)
	}

}

func TestParametersLoading(t *testing.T) {
	filename := "../../secrets/secrets.json"
	address := "vault-staging.betterez.com"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.SkipNow()
	}
	params, err := LoadVaultInfoFromJSONFile(filename, "staging")
	if err != nil {
		t.Fatal(err)
	}
	if params.Address == "" {
		t.Fatal("Bad address loaded")
	}
	if params.Address != address {
		t.Fatalf("expecting %s, got %s", address, params.Address)
	}
}

func TestParametersLoading2(t *testing.T) {
	filename := "../../secrets/secrets.json"
	port := 9000
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.SkipNow()
	}
	params, err := LoadVaultInfoFromJSONFile(filename, "staging")
	if err != nil {
		t.Fatal(err)
	}
	if params.Address == "" {
		t.Fatal("Bad address loaded")
	}
	if params.Port != port {
		t.Fatalf("expecting %d, got %d", port, params.Port)
	}
}
