package main

import (
	"github.com/hhwe/webgo"
	"log"
	"net/http"
	"os"
	"time"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing FindUsers")
	r.Context()
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
	log.Panicln("你好")
	w.Write([]byte("OK"))
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	app := webgo.NewApp(webgo.Recovery, webgo.Logging, dbSession)
	var logger = webgo.Logger

	app.AddRoute("/", http.HandlerFunc(HomePage))
	app.AddRoute("/users", http.HandlerFunc(FindUsers))
	app.AddRoute("/users/:id", http.HandlerFunc(FindUsersByID))
	app.AddRoute("/error", http.HandlerFunc(FindUsersError))

	logger.Info(`
* Environment: production
	WARNING: Do not use the development server in a production environment.
* Debug mode: off
* Running on http://127.0.0.1:8000/ (Press CTRL+C to quit)
`)
	logger.Info(http.ListenAndServe(":"+port, app))
}
