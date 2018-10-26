package main

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// responseWithJson replies to the request with the specified message and HTTP code.
// It does not otherwise end the request; the caller should ensure no further
// writes are done to w.
func responseWithJson(w http.ResponseWriter, code int, msg string, errorCode int, data interface{}) {
	response := Response{
		Code: errorCode,
		Msg:  msg,
		Data: data,
	}
	payload, err := json.Marshal(response)
	if err != nil {
		// logger.Panicln(err)
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(payload)
}

// Jsonify success to response with marshaled json data and 200 http code.
func Jsonify(w http.ResponseWriter, data interface{}) {
	responseWithJson(w, http.StatusOK, "", 0, data)
}

// Abort prevents pending handlers from being called.
func Abort(w http.ResponseWriter, msg string, code int) {
	// logger.Println(msg)
	responseWithJson(w, code, msg, -1, nil)
}

type Resource interface {
	Get()
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

func (u *User) Get(params map[string]string) {
	json.Marshal(u)
}
