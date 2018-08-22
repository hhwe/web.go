package main

import (
	"log"
	"net/http"
	_ "webgo/handles"
)

func httpServer() {
	log.Fatal(http.ListenAndServe(":8888", nil))
}
