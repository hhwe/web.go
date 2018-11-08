package main

import (
	"log"
	"net/http"
	"os"

	"github.com/hhwe/webgo"
	"gopkg.in/mgo.v2"
)

var (
	logger = webgo.Logger
)

func init() {
	logger.Logger.SetFlags(log.LstdFlags)
	logger.SetLever("debug")
}

func main() {
	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	var db, _ = mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{"localhost"},
		Database: "web",
	})
	defer db.Close()

	app := webgo.NewApp(webgo.Recovery, webgo.Logging, dbSession(db))

	app.StaticWeb("/public", ".")
	app.AddRoute("/", http.HandlerFunc(HomePage))
	app.AddRoute("/users", http.HandlerFunc((&User{}).Get))

	app.Run(":8080")
}
