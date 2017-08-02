package btrzaws

import "testing"

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
