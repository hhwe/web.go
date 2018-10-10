package main

import (
	"net/http"
	"time"
)

// session stored in redis
type Session interface {
	Set(key, value interface{}) error // set session value
	Get(key interface{}) interface{}  // get session value
	Delete(key interface{}) error     // delete session value
	SessionID() string                // back current sessionID
}

// addCookie will apply a new cookie to the response of a http
// request, with the key/value this method is passed.
func addCookie(w http.ResponseWriter, name string, value string) {
	expire := time.Now().AddDate(0, 0, 3)
	cookie := http.Cookie{
		Name:    name,
		Value:   value,
		Expires: expire,
	}
	http.SetCookie(w, &cookie)
}


// validate cookie from request of current user
func validCookie(r *http.Request, name string, value string) bool {
	cookie, err := r.Cookie(name)
	if err != nil {
		panic(err)
	}
	if cookie.Value == value {
		return true
	}
	return false
}
