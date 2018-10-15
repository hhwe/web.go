package main

import (
	"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
	ID       bson.ObjectId `bson:"_id" json:"id"`
	Age      int           `bson:"age" json:"age"`
	Sex      int           `bson:"sex" json:"sex"`
	Email    string        `bson:"email" json:"email"`
	Phone    string        `bson:"phone" json:"phone"`
	Summary  string        `bson:"summary"`
	UserName string        `bson:"username"`
	PassWord string        `bson:"password"`
	Created  time.Time     `bson:"created"`
}

func register(w http.ResponseWriter, r *http.Request) {
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
	if err := Insert(db, "user", &u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ResponseWithJson(w, http.StatusCreated, u)
}

func login(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "database").(*mgo.Session)

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := FindOne(db, "user", bson.M{"username": u.UserName,
			"password": u.PassWord}, bson.M{}, &u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// note: if this function execute after json encode will not work
	addCookie(w, "name", u.UserName)

	//if err := json.NewEncoder(w).Encode(user); err != nil {
	//	http.Error(w, err.Error(), http.StatusBadRequest)
	//	return
	//}

	ResponseWithJson(w, http.StatusCreated, u)
}

func selectUser(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "database").(*mgo.Session)

	var users []*User
	if err := db.DB("web").C("user").
		Find(nil).Sort("-created").Limit(100).All(&users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ResponseWithJson(w, http.StatusOK, users)
}

func insertUser(w http.ResponseWriter, r *http.Request) {
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
	if err := Insert(db, "user", &u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ResponseWithJson(w, http.StatusCreated, u)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
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
	if err := FindOne(db, "user", bson.IsObjectIdHex()); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
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
}
