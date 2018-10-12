package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"net/http"
	"time"
	"encoding/json"
)

// TODO: Object Relationship Mapping
type Model interface {
	//GetID(db *mgo.Session) bson.ObjectId
	SelectOne(id interface{}, db *mgo.Session) interface{}  // get one documents
	InsertOne(r io.ReadCloser, db *mgo.Session) error // set one documents
	//UpdateOne(id interface{}) interface{}  // update one documents
	//DeleteOne(id interface{}) error     // delete one documents

	Collection() string  // return collection to storage
}

type BaseModel struct {
	ID bson.ObjectId
	m Model
}

func SelectOne(w http.ResponseWriter, r *http.Request, b *BaseModel) {
	db := context.Get(r, "database").(*mgo.Session)
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		return
	}
	resp := b.m.SelectOne(b.ID, db)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func InsertOne(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "database").(*mgo.Session)

	var users []*User
	if err := db.DB("web").C("user").
		Find(nil).Sort("-created").Limit(100).All(&users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (user *User) SelectOne(id interface{}, db *mgo.Session) (u *User, err error) {
	if err = db.DB("web").C("user").Find(
		bson.M{"username":user.UserName, "password":user.PassWord}).
		Sort("-created").One(&user); err != nil {
		return
	}
	u = user
	return
}

func (user *User) InsertOne(r io.ReadCloser, db *mgo.Session) (err error) {
	if err = json.NewDecoder(r).Decode(&user); err != nil {
		return
	}

	user.ID = bson.NewObjectId()
	user.Created = time.Now()
	if err = db.DB("web").C("user").Insert(&user); err != nil {
		return
	}
	return
}
