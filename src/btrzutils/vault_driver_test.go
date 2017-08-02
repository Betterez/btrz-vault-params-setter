package btrzutils

import (
	"fmt"
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
	if driver.GetVaultStatus() != "online" {
		t.Fatalf("Driver status is %s", driver.GetVaultStatus())
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
	if driver.GetVaultStatus() != "online" {
		t.Fatalf("Driver status is %s", driver.GetVaultStatus())
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
	if driver.GetVaultStatus() != "online" {
		t.Fatalf("Driver status is %s", driver.GetVaultStatus())
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

func TestJSONKeys(t *testing.T) {
	testString := `{"user":"test","password":"q1w2e3"}`
	var results string
	testJSON, err := simplejson.NewJson([]byte(testString))
	if err != nil {
		t.Fatal(err)
	}
	testMap, err := testJSON.Map()
	if err != nil {
		t.Fatal(err, testJSON)
	}
	for key := range testMap {
		keyValue, _ := testJSON.Get(key).String()
		results += keyValue
	}
	if results != "testq1w2e3" {
		t.Fatalf("results = %s", results)
	}
}

func TestJSONAddition(t *testing.T) {
	testString := `{"user":"test","password":"q1w2e3"}`
	additionString := `{"data":"hello world"}`
	resultJSON := simplejson.New()
	additionJSON, err := simplejson.NewJson([]byte(additionString))
	if err != nil {
		t.Fatal(err)
	}
	testJSON, err := simplejson.NewJson([]byte(testString))
	if err != nil {
		t.Fatal(err)
	}
	testMap, err := testJSON.Map()
	if err != nil {
		t.Fatal(err, testJSON)
	}
	for key := range testMap {
		keyValue, _ := testJSON.Get(key).String()
		resultJSON.Set(key, keyValue)
	}
	testMap, err = additionJSON.Map()
	if err != nil {
		t.Fatal(err, testJSON)
	}
	for key := range testMap {
		keyValue, _ := additionJSON.Get(key).String()
		fmt.Printf("%s = %s", key, keyValue)
		resultJSON.Set(key, keyValue)
	}
	resultData, err := resultJSON.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	resutString := string(resultData)
	// note that the list is now ordered
	if resutString != `{"data":"hello world","password":"q1w2e3","user":"test"}` {
		t.Fatalf("bad results: %s\n%v", resutString, additionJSON)
	}

}

func TestExistingJSONValue(t *testing.T) {
	testString := `{"data":{"user":"test","password":"q1w2e3"},"checkpoint":1}`
	JSONData, err := simplejson.NewJson([]byte(testString))
	if err != nil {
		t.Fatal(err)
	}
	checkPoint, err := JSONData.Get("checkpoint").Int()
	if err != nil {
		t.Fatal(err)
	}
	if checkPoint != 1 {
		t.Fatalf("Bad checkpoint value (%d!=1)", checkPoint)
	}
	dataJSON, exists := JSONData.CheckGet("data")
	if !exists {
		t.Fatal("Existing key returned null!")
	}
	userData, err := dataJSON.Get("user").String()
	if err != nil {
		t.Fatal(err)
	}
	if userData != "test" {
		t.Fatalf("Bad value returned, expecting 'test', got '%s'", userData)
	}
}

func TestSetValue(t *testing.T) {
	t.SkipNow()
	token, err := GetToken()
	if err != nil {
		t.Fatal(err)
	}
	if token == "" {
		t.SkipNow()
	}
	connection, err := CreateVaultConnectionFromParameters("vault-staging.betterez.com", token, 9000)
	if err != nil {
		t.Fatal(err)
	}
	randomePath := "secret/" + RandStringRunes(30)
	randomValue := RandStringRunes(30)
	code, err := connection.PutJSONValue(randomePath, fmt.Sprintf(`{"value":"%s"}`, randomValue))
	if err != nil {
		t.Fatal(err)
	}
	if code >= 400 {
		t.Fatalf("bad http code returned - %d", code)
	}
	value, err := connection.GetJSONValue(randomePath)
	if err != nil {
		t.Fatal(err)
	}
	formattedData, err := simplejson.NewJson([]byte(value))
	if err != nil {
		t.Fatal(err)
	}
	fetchedValue, err := formattedData.Get("data").Get("value").String()
	if err != nil {
		t.Fatal(err)
	}
	if fetchedValue != randomValue {
		t.Fatalf("Bad value returned, expecting %s, got %s", randomValue, fetchedValue)
	}
}

func TestAddValueToPath(t *testing.T) {
	t.SkipNow()
	token, err := GetToken()
	randomePath := "secret/" + RandStringRunes(30)
	randomValue := RandStringRunes(30)
	randomValue2 := RandStringRunes(30)
	if err != nil {
		t.Fatal(err)
	}
	if token == "" {
		t.SkipNow()
	}
	connection, err := CreateVaultConnectionFromParameters("vault-staging.betterez.com", token, 9000)
	if err != nil {
		t.Fatal(err)
	}
	code, err := connection.PutJSONValue(randomePath, fmt.Sprintf(`{"value":"%s"}`, randomValue))
	if err != nil {
		t.Fatal(err)
	}
	if code >= 400 {
		t.Fatalf("bad http code returned - %d", code)
	}
	connection.AddValuesInPath(randomePath, fmt.Sprintf(`{"value2":"%s"}`, randomValue2))
	value, err := connection.GetJSONValue(randomePath)
	if err != nil {
		t.Fatal(err)
	}
	formattedData, err := simplejson.NewJson([]byte(value))
	if err != nil {
		t.Fatal(err)
	}
	fetchedValue, err := formattedData.Get("data").Get("value").String()
	if err != nil {
		t.Fatal(err)
	}
	if fetchedValue != randomValue {
		t.Fatalf("Bad value returned, expecting %s, got %s", randomValue, fetchedValue)
	}
	fetchedValue, err = formattedData.Get("data").Get("value2").String()
	if err != nil {
		t.Fatal(err)
	}
	if fetchedValue != randomValue2 {
		t.Fatalf("Bad value returned, expecting %s, got %s", randomValue, fetchedValue)
	}
}
