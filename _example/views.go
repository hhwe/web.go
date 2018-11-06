package main

import (
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
