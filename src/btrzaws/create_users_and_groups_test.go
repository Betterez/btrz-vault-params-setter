package btrzaws

import (
	"encoding/json"
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
