package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing FindUsers")
	w.Write([]byte("Home Page"))
}

func FindUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing FindUsers")
	time.Sleep(time.Second * 1)
	w.Write([]byte("OK"))
}

func FindUsersByID(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing FindUsersByID")
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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	app := NewApp(Recovery, Logging)

	app.AddRoute("/", http.HandlerFunc(HomePage))
	app.AddRoute("/users", http.HandlerFunc(FindUsers))
	app.AddRoute("/users/:id", http.HandlerFunc(FindUsersByID))
	app.AddRoute("/error", http.HandlerFunc(FindUsersError))

	log.Fatal(http.ListenAndServe(":"+port, app))
}
