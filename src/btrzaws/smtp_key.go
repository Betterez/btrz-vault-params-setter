package btrzaws

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

const (
	awsMessageKey          = "SendRawEmail"
	awsMessageVersion byte = 0x02
)

// GenerateSMTPPasswordFromSecret - generate smtp password from a given aws secret
func GenerateSMTPPasswordFromSecret(secret string) (string, error) {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(awsMessageKey))
	value1 := mac.Sum(nil)
	infoWithSignature := append([]byte{}, awsMessageVersion)
	infoWithSignature = append(infoWithSignature, value1...)

	return base64.StdEncoding.EncodeToString(infoWithSignature), nil
}
