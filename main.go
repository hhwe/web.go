package main

import (
	"log"
	"net/http"
	"time"
)

func FindUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing finalHandler")
	time.Sleep(time.Second * 1)
	w.Write([]byte("OK"))
}

func FindUsersByID(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing finalHandler")
	vs := r.URL.Query()
	log.Println(r.URL.RawQuery)
	w.Write([]byte(vs.Get("id")))
}

func FindUsersError(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing finalHandler")
	Logger.Panicln("你好")
	w.Write([]byte("OK"))
}

func main() {
	app := NewApp(Recovery, Logging)
	app.AddRoute("/users", http.HandlerFunc(FindUsers))
	app.AddRoute("/users/:id", http.HandlerFunc(FindUsersByID))
	app.AddRoute("/error", http.HandlerFunc(FindUsersError))

	log.Fatal(http.ListenAndServe(":8080", app))
}
