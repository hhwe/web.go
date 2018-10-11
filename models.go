package main

import (
	"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	"log"
	"net/http"
	"regexp"
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
	Phone    string           `bson:"phone"`
	Summary  string        `bson:"summary"`
	UserName string        `bson:"name"`
	PassWord string        `bson:"password"`
	Created  time.Time     `bson:"created"`
}

func signIn(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "database").(*mgo.Session)

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if u.Phone == "" || u.UserName == "" || u.PassWord == "" {
		http.Error(w, "ERROR: Invalid input", http.StatusBadRequest)
		//logger.Panic("ERROR: Invalid input")
		return
	}
	if m, _ := regexp.MatchString(`^(1[3|4|5|8][0-9]\d{4,8})$`, u.Phone); !m {
		logger.Panic("ERROR: invalid phone")
	}

	u.ID = bson.NewObjectId()
	u.Created = time.Now()
	if err := db.DB("web").C("user").Insert(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	addCookie(w, "name", u.UserName)
}

func auth(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "database").(*mgo.Session)

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.DB("web").C("user").Find(
		bson.M{"username":user.UserName, "password":user.PassWord}).
		Sort("-created").One(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func addUser(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "database").(*mgo.Session)

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if u.Phone == "" || u.UserName == "" || u.PassWord == "" {
		log.Panic("ERROR: need phone and password")
	}

	u.ID = bson.NewObjectId()
	u.Created = time.Now()
	if err := db.DB("web").C("user").Insert(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//http.Redirect(w, r, "/login", http.StatusFound)
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
