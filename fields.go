package main

import (
	"encoding/json"
	"net/http"
)

// todo: response error code table
const (
	StatusFail = -1
	StatusSuccess = 0
)

var responseCode = map[int]string{
	StatusFail:           "Fail",
	StatusSuccess: "Success",
}

func ResponseCodeText(code int) string {
	return responseCode[code]
}

type Response struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}

func ResponseWithJson(w http.ResponseWriter, code int, msg string, data interface{}) {
	response := Response{
		Code:code,
		Msg:msg,
		Data:data,
	}
	payload, err := json.Marshal(response)
	if err != nil {
		logger.Panicln(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(payload)
}
