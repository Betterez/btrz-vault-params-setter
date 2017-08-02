package btrzaws

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/service/iam"
)

func createTestPolicy() {
	policyDocument := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":["ec2.amazonaws.com"]},"Action":["sts:AssumeRole"]}]}`
	policyName := "vue-live-test"
	//userName := "testUser"
	path := "/"
	policyDescription := "policy for vue test service"
	awsSession, err := GetAWSSession()
	if err != nil {
		fmt.Print(err, "can't get a session")
		os.Exit(1)
	}
	log.Println("session created")
	iamService := iam.New(awsSession)
	if iamService == nil {
		fmt.Println("can't create iam")
		os.Exit(1)
	}
	policyResponse, err := iamService.CreatePolicy(
		&iam.CreatePolicyInput{
			PolicyDocument: &policyDocument,
			Path:           &path,
			Description:    &policyDescription,
			PolicyName:     &policyName,
		})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(policyResponse)
}