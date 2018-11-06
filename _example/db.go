package main

import (
	"context"
	"gopkg.in/mgo.v2"
	"net/http"
)

var db, _ = mgo.DialWithInfo(&mgo.DialInfo{
	Addrs:    []string{"localhost"},
	Database: "web",
})

func dbSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := db.Copy()
		defer session.Close()

		ctx := context.WithValue(r.Context(), "database", session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
