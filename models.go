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
	db := context.Get(r, "database").(*mgo.Database)

	// decode the request body
	var c comment
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// give the comment a unique ID and set the time
	c.ID = bson.NewObjectId()
	c.When = time.Now()
	// insert it into the database
	if err := db.C("comments").Insert(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//// redirect to it
	//http.Redirect(w, r, "/comments/"+c.ID.Hex(), http.StatusTemporaryRedirect)
}

func handleRead(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "database").(*mgo.Database)
	// load the comments
	var comments []*comment
	if err := db.C("comments").Find(nil).Sort("-when").
			Limit(100).All(&comments); err != nil {
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
	Summary  string        `bson:"summary" json:"summary"`
	UserName string        `bson:"username" json:"username"`
	PassWord string        `bson:"password" json:"password"`
	Created  time.Time     `bson:"created" json:"created"`
}

func register(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "database").(*mgo.Database)

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if u.Phone == "" || u.UserName == "" || u.PassWord == "" {
		http.Error(w, "ERROR: Invalid input", http.StatusBadRequest)
		//logger.Panic("ERROR: Invalid input")
		return
	}
	if m, _ := regexp.MatchString(`^(1[3458][0-9]\d{4,8})$`, u.Phone); !m {
		logger.Panic("ERROR: invalid phone")
	}

	u.ID = bson.NewObjectId()
	u.Created = time.Now()
	// Prevent plaintext Passwords are stored
	u.PassWord = HashSha256(u.PassWord)
	if err := db.C("user").Insert(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ResponseWithJson(w, http.StatusCreated, ResponseCodeText(StatusSuccess), nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "database").(*mgo.Database)

	var u User
	// todo: verify the validity of data, csrf token
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	u.PassWord = HashSha256(u.PassWord)
	if err := db.C("user").Find(bson.M{"username": u.UserName,
			"password": u.PassWord}).One(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token:=addToken()
	// note: if this function execute after json encode will not work
	addCookie(w, "name", u.UserName)

	//if err := json.NewEncoder(w).Encode(user); err != nil {
	//	http.Error(w, err.Error(), http.StatusBadRequest)
	//	return
	//}

	ResponseWithJson(w, http.StatusOK, ResponseCodeText(StatusSuccess), token)
}

func findUser(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "database").(*mgo.Database)

	var users []*User
	if err := db.C("user").Find(nil).Sort("-created").
			Limit(100).All(&users); err != nil {
		logger.Panicln(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ResponseWithJson(w, http.StatusOK, "", users)
}

func insertUser(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "database").(*mgo.Database)

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if u.Phone == "" || u.UserName == "" || u.PassWord == "" {
		log.Panic("ERROR: need phone and password")
	}

	u.ID = bson.NewObjectId()
	u.Created = time.Now()
	// Prevent plaintext Passwords are stored
	u.PassWord = HashSha256(u.PassWord)
	if err := db.C("user").Insert(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ResponseWithJson(w, http.StatusCreated, "", u)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "database").(*mgo.Database)

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if u.Phone == "" || u.UserName == "" || u.PassWord == "" {
		log.Panic("ERROR: need phone and password")
	}
	// Prevent plaintext Passwords are stored
	u.PassWord = HashSha256(u.PassWord)
	if err := db.C("user").Find(bson.M{"id": u.ID}).One(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ResponseWithJson(w, http.StatusOK, "", u)
}

func removeUser(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "database").(*mgo.Database)

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if u.Phone == "" || u.UserName == "" || u.PassWord == "" {
		log.Panic("ERROR: need phone and password")
	}

	u.ID = bson.NewObjectId()
	u.Created = time.Now()
	if err := db.C("user").RemoveId(u.ID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ResponseWithJson(w, http.StatusOK, "", u)
}
