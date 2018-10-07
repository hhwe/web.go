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
	if err := db.DB("commentsapp").C("comments").Insert(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//// redirect to it
	//http.Redirect(w, r, "/comments/"+c.ID.Hex(), http.StatusTemporaryRedirect)
}

func handleRead(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r,"database").(*mgo.Session)
	// load the comments
	var comments []*comment
	if err := db.DB("commentsapp").C("comments").
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

////个人项目部分代码
//type User struct {
//	ID       bson.ObjectId `bson:"_id"`
//	UserName string        `bson:"username"`
//	Summary  string        `bson:"summary"`
//	Age      int           `bson:"age"`
//	Phone    int           `bson:"phone"`
//	PassWord string        `bson:"password"`
//	Sex      int           `bson:"sex"`
//	Name     string        `bson:"name"`
//	Email    string        `bson:"email"`
//}
//
//func AddUser(password string, username string) (err error) {
//	con := GetDataBase().C("user")
//	//可以添加一个或多个文档
//	/* 对应mongo命令行
//	   db.user.insert({username:"13888888888",summary:"code",
//	   age:20,phone:"13888888888"})*/
//	err = con.Insert(&User{ID: bson.NewObjectId(), UserName: username, PassWord: password})
//	return
//}
//
//func FindUser(username string) (User, error) {
//	var user User
//	con := GetDataBase().C("user")
//	//通过bson.M(是一个map[string]interface{}类型)进行
//	//条件筛选，达到文档查询的目的
//	/* 对应mongo命令行
//	  db.user.find({username:"13888888888"})*/
//	if err := con.Find(bson.M{"username": username}).One(&user); err != nil {
//		if err.Error() != GetErrNotFound().Error() {
//			return user, err
//		}
//
//	}
//	return user, nil
//}
