// This is the viewpoint of web application.

// At the beginning, we use templates to parse all static file,
// so that we don't need parse them every request.
// The renderTemplate functions render a template file when get some page.

// Every request handled depend on their's method,
// different method has there own's handle function.
// The low-level handle functions is in models.go.
package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

// parse static file when app was compiled
var templates = template.Must(template.ParseFiles(
	"static/index.html", "static/register.html", "static/login.html"))

// logging setting
var logger = log.New(os.Stderr, "", log.Ldate|log.Ltime)

func renderTemplate(w http.ResponseWriter, tmpl string) {
	csrf := addToken("csrf")
	err := templates.ExecuteTemplate(w, tmpl+".html", csrf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index")
}

func Register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		renderTemplate(w, r.URL.Path[1:])
	case "POST":
		register(w, r)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		renderTemplate(w, r.URL.Path[1:])
	case "POST":
		login(w, r)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index")
}

func Cart(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, r.URL.Path[1:])
}

func Comment(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleRead(w, r)
	case "POST":
		handleInsert(w, r)
	}
}

func ErrorTest(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func UserApi(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		findUser(w, r)
	case "POST":
		insertUser(w, r)
	case "PUT":
		updateUser(w, r)
	case "DELETE":
		removeUser(w, r)
	}
}
