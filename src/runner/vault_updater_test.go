package main

import (
	"encoding/json"
	"os"
	"testing"
)

func TestLoadingValues(t *testing.T) {
	const (
		filename       = "../../test_objects/vault_params.json"
		expectedServer = "vault.myserver.com"
	)
	var err error
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		t.Skipf("no file %s", filename)
	}
	info, err := loadVaultInfoFromFile(filename)
	if err != nil {
		t.Error(err)
	}
	if info.Server.Address != expectedServer {
		t.Fatalf("expected %s, got %s", expectedServer, info.Server.Address)
	}
	if info.Values["taltul"] == nil {
		t.Fatal("bad mapping")
	}
	if info.Values["taltul"]["pass"] != "some" {
		t.Fatalf("bad map value: wanted %s, got %s", "some", info.Values["taltul"]["pass"])
	}
	jsonData, err := json.Marshal(info.Values["taltul"])
	if err != nil {
		t.Fatalf("error %v marshaling", err)
	}
	jsonString := string(jsonData)
	if jsonString != `{"data":"123","pass":"some"}` {
		t.Fatalf("bad json string: %s", jsonString)
	}
}
