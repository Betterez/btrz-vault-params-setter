package btrzdb

import (
	"os"
	"testing"
)

func TestDatabasaParameters(t *testing.T) {
	data, err := GetDialInfo("local1")
	if os.IsNotExist(err) {
		t.SkipNow()
	}
	if err != nil {
		t.Fatal(err)
	}
	if data.Addrs[0] != "192.168.0.41" {
		t.Fatal("bad server address")
	}
	if data.Username != "tal" {
		t.Fatal("bad username")
	}
}
