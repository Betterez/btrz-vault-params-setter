package btrzaws

import "testing"

const (
	testKey = "q1w2e3r4t5y6q1w2e3r4t5y6q1w2e3r4t5y6zzzz"
	result  = "Alo/zZmiD2A9FiTR2Y1hitH3JZlcHOf6jTIwFBPlJw6b"
)

func TestGenerator(t *testing.T) {
	value, err := GenerateSMTPPasswordFromSecret(testKey)
	if err != nil {
		t.Fatal("error converting:", err)
	}
	if value != result {
		t.Fatalf("values don't match!\nexpected %s\ngot %s", result, value)
	}
}
