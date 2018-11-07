package main

import (
	"fmt"
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

	app.AddRoute("/", http.HandlerFunc(HomePage))

	logger.Info(fmt.Sprintf(`webgo restful api application 
* Environment: production
* Debug mode: off
* Running on http://%s:%s/ (Press CTRL+C to quit)
`, host, port))
	logger.Info(http.ListenAndServe(":"+port, app))
}
