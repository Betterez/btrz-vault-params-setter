package btrzaws

import (
	"encoding/json"
	simplejson "github.com/bitly/go-simplejson"
	"io/ioutil"
	"testing"
)

func TestBadServiceInfo(t *testing.T) {
	info := &ServiceInformation{}
	if info.IsInformationOK() {
		t.Fatal("Bad service info data passing as ok")
	}
}

func TestBadServiceInfo2(t *testing.T) {
	info := GenerateServiceInformation("test")
	if info.IsInformationOK() {
		t.Fatal("Bad service info data passing as ok")
	}
}

func TestAddArn(t *testing.T) {
	info := GenerateServiceInformation("test")
	info.AddServiceArn("test")
	if !info.IsInformationOK() {
		t.Fatal("Service information is correct but bad value returned")
	}
}
func TestLoadingServiceFile(t *testing.T) {
	servicesFile := "../../services/services.json"
	data, err := ioutil.ReadFile(servicesFile)
	if err != nil {
		t.Fatal(err)
	}
	servicesData := make(map[string]ServiceInformation)
	err = json.Unmarshal(data, &servicesData)
	if err != nil {
		t.Fatal(err)
	}
	for serviceKey := range servicesData {
		serviceInfo := servicesData[serviceKey]
		if serviceInfo.ServiceName != serviceKey {
			t.Fatalf("service name error, expecting btrz-data-import, got %s", serviceInfo.ServiceName)
		}
	}
}

func TestMongoInformationDatabaseName(t *testing.T) {
	servicesFile := "../../services/services.json"
	data, err := ioutil.ReadFile(servicesFile)
	if err != nil {
		t.Fatal(err)
	}
	servicesData := make(map[string]ServiceInformation)
	err = json.Unmarshal(data, &servicesData)
	if err != nil {
		t.Fatal(err)
	}
	information, ok := servicesData["btrz-worker-data-import"]
	if !ok {
		t.Fatal("failed to get value")
	}
	if information.MongoSettings.DatabaseName["staging"] != "bz_staging" {
		t.Fatal("Can't retrieve staging db info")
	}

}
func TestMongoInformationRoleName(t *testing.T) {
	servicesFile := "../../services/services.json"
	expectedRole := "dbOwner"
	data, err := ioutil.ReadFile(servicesFile)
	if err != nil {
		t.Fatal(err)
	}
	servicesData := make(map[string]ServiceInformation)
	err = json.Unmarshal(data, &servicesData)
	if err != nil {
		t.Fatal(err)
	}
	information, ok := servicesData["btrz-worker-data-import"]
	if !ok {
		t.Fatal("failed to get value")
	}
	if information.MongoSettings.DatabaseRole != expectedRole {
		t.Fatalf("Bad database role. expecting %s, got %s", expectedRole, information.MongoSettings.DatabaseRole)
	}
}

func TestServiceInfoWithoutMongo(t *testing.T) {
	t.SkipNow()
	servicesFile := "../../services/services.json"
	data, err := ioutil.ReadFile(servicesFile)
	if err != nil {
		t.Fatal(err)
	}
	servicesData := make(map[string]ServiceInformation)
	err = json.Unmarshal(data, &servicesData)
	if err != nil {
		t.Fatal(err)
	}
	information, ok := servicesData["btrz-worker-exports"]
	if !ok {
		t.Fatal("failed to get value")
	}
	if information.MongoSettings.DatabaseRole != "" {
		t.Fatal("got a value for non existing value")
	}
}

func TestJSONStringCreation(t *testing.T) {
	tester := make(map[string]string)
	tester["one"] = "111"
	tester["two"] = "222"
	js, err := simplejson.NewJson([]byte(createJSONStringFromKeyValues(tester)))
	if err != nil {
		t.Fatalf("%v error getting json string", err)
	}
	value, err := js.Get("one").String()
	if value != "111" {
		t.Fatalf("Wrong values received %s ", value)
	}
}
