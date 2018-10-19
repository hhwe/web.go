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
	insertUser(w, r)
}

// todo: add to config file
var loginCookieName = "_uid"

func login(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "database").(*mgo.Database)

	var u User
	// todo: verify the validity of data, csrf token
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		Abort(w, "error request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	u.PassWord = HashSHA256(u.PassWord)
	if err := db.C("user").Find(bson.M{"username": u.UserName}).One(&u); err != nil {
		Abort(w, "user not fount", http.StatusBadRequest)
		return
	}

	token := GenerateToken(string(u.ID))
	// note: if this function execute after json encode will not work
	addCookie(w, loginCookieName, string(u.ID))

	Jsonify(w, token)
}

func findUser(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "database").(*mgo.Database)

	var users []*User
	if err := db.C("user").Find(nil).Sort("-created").
		Limit(100).All(&users); err != nil {
		logger.Panicln(err)
		Abort(w, "error request body", http.StatusBadRequest)
		return
	}

	Jsonify(w, users)
}

func insertUser(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "database").(*mgo.Database)

	var u User
	// todo: prevent XSS can use http.Eascape().
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		Abort(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if m, _ := regexp.MatchString(`^[\p{Han}\w]{2,8}$`, u.UserName); !m {
		Abort(w, "invalid name, only chainese or alpnum are permitted", http.StatusBadRequest)
		return
	}

	if m, _ := regexp.MatchString(`^(1[3458][0-9]\d{4,8})$`, u.Phone); !m {
		Abort(w, "invalid phone", http.StatusBadRequest)
		return
	}

	if m, _ := regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, u.Email); !m {
		Abort(w, "invalid email", http.StatusBadRequest)
		return
	}

	if m, _ := regexp.MatchString(`^[\w._]{6,20}$`, u.PassWord); !m {
		Abort(w, "invalid password", http.StatusBadRequest)
		return
	}

	u.ID = bson.NewObjectId()
	u.Created = time.Now()
	// prevent plaintext passwords are stored
	u.PassWord = HashSHA256(u.PassWord)
	// ensure user field is unique
	au := User{}
	db.C("user").Find(bson.M{"$or": []bson.M{
		bson.M{"username": u.UserName},
		bson.M{"email": u.Email},
		bson.M{"phone": u.Phone},
	}}).One(&au) // ignore error when query data
	if au.ID != "" {
		Abort(w, "repeated phone, email or username", http.StatusBadRequest)
		return
	}

	if err := db.C("user").Insert(&u); err != nil {
		panic(err)
	}

	Jsonify(w, nil)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "database").(*mgo.Database)

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		Abort(w, "error request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if u.Phone == "" || u.UserName == "" || u.PassWord == "" {
		log.Panic("ERROR: need phone and password")
	}
	// Prevent plaintext Passwords are stored
	u.PassWord = HashSHA256(u.PassWord)
	if err := db.C("user").Find(bson.M{"id": u.ID}).One(&u); err != nil {
		Abort(w, "error request body", http.StatusBadRequest)
		return
	}

	Jsonify(w, u)
}

func removeUser(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "database").(*mgo.Database)

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		Abort(w, "error request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if u.Phone == "" || u.UserName == "" || u.PassWord == "" {
		log.Panic("ERROR: need phone and password")
	}

	u.ID = bson.NewObjectId()
	u.Created = time.Now()
	if err := db.C("user").RemoveId(u.ID); err != nil {
		Abort(w, "error request body", http.StatusBadRequest)
		return
	}

	Jsonify(w, u)
}
