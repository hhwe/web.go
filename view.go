package main

import (
	"log"
	"net/http"
	"os"
	"html/template"
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

