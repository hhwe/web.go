package main

import (
	"log"
	"net/http"
	_ "webgo/views"
)

func main() {
	log.Fatal(http.ListenAndServe(":8888", nil))
}
