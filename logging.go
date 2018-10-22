package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"time"
)

type Logger struct {

}

func (l Logger) ServerHttp(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	next(w, r)

}

func NewLogger()  {
	log.New(os.Stderr, "", log.Ldate|log.Ltime)

}
