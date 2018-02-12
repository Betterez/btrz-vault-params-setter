package btrzutils

import "testing"

func TestRandomizerStringLength(t *testing.T) {
	result := RandStringRunes(30)
	if len(result) != 30 {
		t.Fatal("Bad string length returned")
	}
}

func TestRandomString(t *testing.T) {
	str1 := RandStringRunes(30)
	str2 := RandStringRunes(30)
	if str1 == str2 {
		t.Fatalf("Strings are the same %s==%s", str1, str2)
	}

}
