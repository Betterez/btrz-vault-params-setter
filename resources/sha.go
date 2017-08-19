package btrztest

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"
)

func createSha(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}
func main1() {
	fmt.Println("")
	ref := time.Now().UTC()
	fmt.Println(ref.Format("Mon, _2 Jan 2006 15:04:05 GMT"))
	//fmt.Println(ref.Format(time.UTC))
}
