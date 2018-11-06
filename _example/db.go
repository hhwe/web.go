package main

import (
	"context"
	"net/http"

	"github.com/hhwe/webgo"
	"gopkg.in/mgo.v2"
)

func dbSession(db *mgo.Session) webgo.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := db.Copy()
			defer session.Close()

			ctx := context.WithValue(r.Context(), "database", session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
