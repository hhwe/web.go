package main

import (
	"github.com/hhwe/webgo"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"gopkg.in/mgo.v2"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("database").(*mgo.Session)
	u := User{}
	err := db.DB("web").C("user").Find(bson.M{}).One(&u)
	if err != nil {
		panic(err)
	}

	w.Write([]byte(u.Username))
}

type UserController struct {
	model User
	webgo.Resource
}

func (u *UserController) Get(w http.ResponseWriter, r *http.Request) {
	db := r.Context().Value("database").(*mgo.Session)
	defer db.Close()

	err := db.DB("web").C("user").Find(bson.M{}).One(&u.model)
	if err != nil {
		panic(err)
	}
	w.Write(u.model.Marshal())
}

func (u *UserController) Post(w http.ResponseWriter, r *http.Request) {

}
