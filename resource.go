package main

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Resource interface {
	Get()
}

type User struct {
	ID       bson.ObjectId `bson:"_id" json:"id"`
	Age      int           `bson:"age" json:"age"`
	Sex      int           `bson:"sex" json:"sex"`
	Email    string        `bson:"email" json:"email"`
	Phone    string        `bson:"phone" json:"phone"`
	Summary  string        `bson:"summary" json:"summary"`
	UserName string        `bson:"username" json:"username"`
	PassWord string        `bson:"password" json:"password"`
	Created  time.Time     `bson:"created" json:"created"`
}

func (u *User) Get(params map[string]string) {
}
