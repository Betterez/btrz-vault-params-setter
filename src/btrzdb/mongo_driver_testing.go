package btrzdb

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
)

func testRecord(driver *mgo.Session) {
	dbName := "sales"
	dbCollection := "recipts"
	record := make(map[string]interface{})
	record["name"] = "dudu"
	record["partner"] = "gulu"
	record["price"] = 18.25
	record["address"] = map[string]interface{}{"street": "Observatory lane", "number": 23}
	err := driver.DB(dbName).C(dbCollection).Insert(record)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	backData := driver.DB(dbName).C(dbCollection).Find(bson.M{"name": "dudu"})
	looper := backData.Iter()
	result := make(map[string]interface{})
	for looper.Next(result) {
		fmt.Println(result)
	}
}

func listUsers(driver *mgo.Session) error {
	dbCollection := "system.users"
	usersData := driver.DB("admin").C(dbCollection).Find(bson.M{})
	result := make(map[string]interface{})
	userItem := usersData.Iter()
	fmt.Println(" ==============================================")
	fmt.Println("|\tuser\t\t|\tdatabase\t|")
	fmt.Println(" ==============================================")
	for userItem.Next(result) {
		fmt.Printf("|\t%s\t\t|\t%s\t\t|\n", result["user"], result["db"])
	}
	fmt.Println(" ===============================================")
	return nil
}
