package btrzaws

import "fmt"

// AwsUserInfo - manage aws user info from iam
type AwsUserInfo struct {
	Username   string
	Policies   []string
	AccessKeys []string
}

// ToString - format the info into a string
func (info *AwsUserInfo) ToString() string {
	return fmt.Sprintf("Information for %s:\n======================\nPolicies:\n=============\n%v\nAccess keys:\n======================\n%v\n",
		info.Username, info.Policies, info.AccessKeys)
}
