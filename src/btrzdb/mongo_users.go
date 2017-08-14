package btrzdb

import (
	"errors"
	"fmt"
	simplejson "github.com/bitly/go-simplejson"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
)

// CreateUser - create a user for the mentioned environment, unless one exists
func CreateUser(username, password, databaseName, databaseRole, environment string) (bool, error) {
	deployment, err := GetDialInfo(environment)
	if err != nil {
		return false, err
	}
	connection, err := mgo.DialWithInfo(deployment)
	if err != nil {
		return false, err
	}
	dbCollection := "system.users"
	usersData := connection.DB("admin").C(dbCollection).Find(bson.M{})
	result := make(map[string]interface{})
	userItem := usersData.Iter()
	found := false
	for userItem.Next(result) && !found {
		if username == result["user"] || result["db"] == databaseName {
			found = true
		}
	}
	if found {
		return false, nil
	}
	betterezUser := &mgo.User{Password: password, Username: username, Roles: []mgo.Role{mgo.Role(databaseRole)}}
	err = connection.DB(databaseName).UpsertUser(betterezUser)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetDialInfo - get dial info from secret file
func GetDialInfo(environment string) (*mgo.DialInfo, error) {
	const fileName = "./secrets/settings.json"
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return nil, err
	}
	fileHandler, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	result := &mgo.DialInfo{}
	databaseData, err := simplejson.NewFromReader(fileHandler)
	if err != nil {
		return nil, err
	}
	result.Username, err = databaseData.Get(environment).Get("mongo").Get("username").String()
	if err != nil {
		fmt.Println(databaseData)
		return nil, errors.New("can't get username")
	}
	result.Password, err = databaseData.Get(environment).Get("mongo").Get("password").String()
	if err != nil {
		return nil, err
	}
	value, err := databaseData.Get(environment).Get("mongo").Get("address").String()
	if err != nil {
		return nil, err
	}
	result.Addrs = []string{value}
	return result, nil
}
