package webgo

import (
	"encoding/json"
	"net/http"
)

// assert assssment expression is true
func assert(exp bool, text string) {
	if !exp {
		panic(text)
	}
}

type Response struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// responseWithJson replies to the request with the specified message and HTTP code.
// It does not otherwise end the request; the caller should ensure no further
// writes are done to w.
func responseWithJson(w http.ResponseWriter, code int, msg string, data interface{}) {
	response := Response{
		Msg:  msg,
		Data: data,
	}
	payload, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(payload)
}

// Jsonify success to response with marshaled json data and 200 http code.
func Jsonify(w http.ResponseWriter, data interface{}) {
	responseWithJson(w, http.StatusOK, "", data)
}

// Abort prevents pending handlers from being called.
func Abort(w http.ResponseWriter, msg string, code int) {
	responseWithJson(w, code, msg, nil)
}
