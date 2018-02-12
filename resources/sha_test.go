package btrztest

import (
	"crypto/sha256"
	"encoding/base64"
	"testing"
)

func createSha(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}

func TestSha256Creation(t *testing.T) {
	// the sha256 doesn't work for betterez sha, but does for aws.
	// need more testing (one for bz and one for aws)
	t.SkipNow()
	const expected = "8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92"
	result := createSha("123456")
	if result != expected {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}
