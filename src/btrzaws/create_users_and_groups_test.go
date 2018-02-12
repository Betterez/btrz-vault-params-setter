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
	if !info.IsInformationOK() {
		t.Fatal("failed to create general service")
	}
}

func TestAddArn(t *testing.T) {
	info := GenerateServiceInformation("test")
	info.AddServiceArn("test")
	if !info.IsInformationOK() {
		t.Fatal("Service information is correct but bad value returned")
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
	if err != nil {
		t.Fatalf("%v error getting json string", err)
	}
	if value != "111" {
		t.Fatalf("Wrong values received %s ", value)
	}
}
