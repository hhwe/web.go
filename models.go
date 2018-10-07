package main

import (
	"encoding/json"
	"github.com/gorilla/context"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

type comment struct {
	ID     bson.ObjectId `json:"id" bson:"_id"`
	Author string        `json:"author" bson:"author"`
	Text   string        `json:"text" bson:"text"`
	When   time.Time     `json:"when" bson:"when"`
}

func handleInsert(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "database").(*mgo.Session)

	// decode the request body
	var c comment
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// give the comment a unique ID and set the time
	c.ID = bson.NewObjectId()
	c.When = time.Now()
	// insert it into the database
	if err := db.DB("web").C("comments").Insert(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//// redirect to it
	//http.Redirect(w, r, "/comments/"+c.ID.Hex(), http.StatusTemporaryRedirect)
}

func handleRead(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "database").(*mgo.Session)
	// load the comments
	var comments []*comment
	if err := db.DB("web").C("comments").
		Find(nil).Sort("-when").Limit(100).All(&comments); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// write it out
	if err := json.NewEncoder(w).Encode(comments); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type User struct {
	ID       bson.ObjectId `bson:"_id"`
	Age      int           `bson:"age"`
	Sex      int           `bson:"sex"`
	Email    string        `bson:"email"`
	Phone    int           `bson:"phone"`
	Summary  string        `bson:"summary"`
	UserName string        `bson:"name"`
	PassWord string        `bson:"password"`
	Created  time.Time     `bson:"created"`
}

func addUser(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "database").(*mgo.Session)

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u.ID = bson.NewObjectId()
	u.Created = time.Now()
	if err := db.DB("web").C("user").Insert(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
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
