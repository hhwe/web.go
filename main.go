package main

import (
	"gopkg.in/mgo.v2"
	"html/template"
	"log"
	"net/http"
	"os"
)

// parse static file when app was compiled
var templates = template.Must(template.ParseFiles(
	"static/index.html", "static/register.html", "static/login.html"))

// logging setting
var logger = log.New(os.Stderr, "", log.Ldate | log.Ltime | log.Lshortfile)

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", nil)
}

func Register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		renderTemplate(w, r.URL.Path[1:], nil)
	case "POST":
		signIn(w, r)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r.URL.Path[1:], nil)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r.URL.Path[1:], nil)
}

func Cart(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r.URL.Path[1:], nil)
}

func Comment(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleRead(w, r)
	case "POST":
		handleInsert(w, r)
	}
}

func main() {
	// connect to the database
	db, err := mgo.Dial("localhost")
	if err != nil {
		logger.Fatal("cannot dial mongo", err)
	}
	defer db.Close() // clean up when we’re done

	// testing db
	http.HandleFunc("/comments", Chain(Comment, Method("GET", "POST"), DBSession(db), Logging()))

	// home page
	http.HandleFunc("/", Chain(Index, Method("GET"), DBSession(db), Logging()))

	// login and sign in
	http.HandleFunc("/register", Chain(Register, Method("GET", "POST"), DBSession(db), Logging()))
	http.HandleFunc("/login", Chain(Login, Method("GET", "POST"), DBSession(db), Logging()))
	http.HandleFunc("/logout", Chain(Logout, Method("GET", "POST"), DBSession(db), Logging()))

	// shopping cart settlement
	http.HandleFunc("/cart", Chain(Cart, Method("GET", "POST"), Logging()))

	// handle static file with file server
	http.Handle("/static/", http.StripPrefix(
		"/static/", http.FileServer(http.Dir("static/"))))

	logger.Println(` 
* Environment: production
  WARNING: Do not use the development server in a production environment.
* Debug mode: off
* Running on http://127.0.0.1:8000/ (Press CTRL+C to quit)
`)
	logger.Fatal(http.ListenAndServe(":8000", nil))

}
