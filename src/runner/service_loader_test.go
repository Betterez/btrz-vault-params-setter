package main

import (
	"btrzaws"
	"encoding/json"
	"io/ioutil"
	"testing"
)

const (
	servicesFile = "../../sample_files/services.json"
)

func TestServiceLoading(t *testing.T) {
	data, err := ioutil.ReadFile(servicesFile)
	if err != nil {
		t.Fatal(err)
	}
	servicesData := make(map[string]btrzaws.ServiceInformation)
	err = json.Unmarshal(data, &servicesData)
	if err != nil {
		t.Fatal(err)
	}
}
