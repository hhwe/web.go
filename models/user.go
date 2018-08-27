package models

import (
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID       bson.ObjectId `bson:"_id" json:"id"`
	Name     string        `form:"name" json:"name" xml:"name"  binding:"required"`
	Password string        `form:"password" json:"password" xml:"password" binding:"required"`
}

const (
	userCollection = "users"
)

func (u *User) FindAll() (users []User, err error) {
	err = db.C(userCollection).Find(bson.M{}).All(&users)
	return
}

func (u *User) FindById(id string) (user User, err error) {
	err = db.C(userCollection).Find(bson.ObjectIdHex(id)).One(&user)
	return
}

func (u *User) Insert(user User) (err error) {
	err = db.C(userCollection).Insert(&user)
	return
}

func (u *User) Delete(user User) (err error) {
	err = db.C(userCollection).Remove(&user)
	return
}

func (u *User) Update(user User) (err error) {
	err = db.C(userCollection).UpdateId(user.ID, &user)
	return
}
