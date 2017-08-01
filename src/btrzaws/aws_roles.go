package btrzaws

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/service/iam"
)

func createTestRole() {
	policyDocument := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":["ec2.amazonaws.com"]},"Action":["sts:AssumeRole"]}]}`
	roleName := "TheTestRole2"
	policyName := "arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess"
	userName := "testUser"
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
	path := "/"
	resp, err := iamService.CreateRole(&iam.CreateRoleInput{
		AssumeRolePolicyDocument: &policyDocument,
		RoleName:                 &roleName,
		Path:                     &path,
	})
	if err != nil {
		fmt.Print(err, "exiting", "\n")
		os.Exit(1)
	}
	fmt.Println(resp.String(), "created.")
	policyResponse, err := iamService.AttachRolePolicy(&iam.AttachRolePolicyInput{
		RoleName:  &roleName,
		PolicyArn: &policyName,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(policyResponse)

	userOutput, err := iamService.CreateUser(&iam.CreateUserInput{
		UserName: &userName,
	})
	if err != nil {
		fmt.Print(err, "exiting", "\n")
		os.Exit(1)
	}
	fmt.Println(userOutput)
}
