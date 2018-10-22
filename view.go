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
	"io"
	"net/http"
)

// parse static file when app was compiled
var templates = template.Must(template.ParseFiles(
	"static/index.html", "static/register.html", "static/login.html"))

func renderTemplate(w http.ResponseWriter, tmpl string) {
	csrf := GenerateToken("csrf")
	err := templates.ExecuteTemplate(w, tmpl+".html", csrf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	// Before any call to WriteHeader or Write, declare
	// the trailers you will set during the HTTP
	// response. These three headers are actually sent in
	// the trailer.
	w.Header().Set("Trailer", "AtEnd1, AtEnd2")
	w.Header().Add("Trailer", "AtEnd3")

	w.Header().Set("Content-Type", "text/plain; charset=utf-908") // normal header
	w.WriteHeader(http.StatusOK)

	w.Header().Set("AtEnd1", "value 1")
	io.WriteString(w, "This HTTP response has both headers before this text and trailers at the end.\n")
	w.Header().Set("AtEnd2", "value 2")
	w.Header().Set("AtEnd3", "value 3") // These will appear as trailers.

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
