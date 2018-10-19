package main

import (
	"encoding/json"
	"net/http"
)

// todo: response error code table
const (
	StatusFail    = -1
	StatusSuccess = 0
)

var responseCode = map[int]string{
	StatusFail:    "Fail",
	StatusSuccess: "Success",
}

func ResponseCodeText(code int) string {
	return responseCode[code]
}

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
		logger.Panicln(err)
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
func Abort(w http.ResponseWriter, msg string, code int)  {
	logger.Println(msg)
	responseWithJson(w, code, msg, -1, nil)
}
