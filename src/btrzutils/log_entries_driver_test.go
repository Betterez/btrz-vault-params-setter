package btrzutils

import (
	"os"
	"testing"
)

func TestLogEntriesConnection(t *testing.T) {
	const fileName = "../../secrets/log_entries.json"
	t.SkipNow()
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		t.SkipNow()
	}
	connection, err := CreateConnectionFromSecretsFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	_, err = connection.GetUsers()
	if err != nil {
		t.Fatal(err)
	}
}
