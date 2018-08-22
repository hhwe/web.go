package models

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// { "_id" : ObjectId("5b2317b0a5dbcbde017c692c"), "i" : 0, "username" : "user0", "age" : 17, "created" : ISODate("2018-06-15T01:34:40.675Z") }
type User struct {
	ID       bson.ObjectId       `bson:"_id" json:"id"`
	I        int                 `bson:"i" json:"i"`
	UserName string              `bson:"name" json:"name"`
	Age      int                 `bson:"age" json:"age"`
	Created  bson.MongoTimestamp `bson:"created" json:"created"`
}

func MongoTest() {
	url := "127.0.0.1"
	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}

	database := "test"
	collection := "users"
	c := session.DB(database).C(collection)

	result := User{}
	query := bson.M{"age": 33}
	err = c.Find(query).One(&result)
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}
