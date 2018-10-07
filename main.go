package main

import (
	"html/template"
	"log"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("static/index.html"))
	tmpl.Execute(w, nil)
}

func Register(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("static/register.html"))
	tmpl.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", Chain(Index, Method("GET"), Logging()))

	http.HandleFunc("/register", Chain(Register, Method("GET"), Logging()))

	// handle static file with file server
	http.Handle("/static/", http.StripPrefix(
		"/static/", http.FileServer(http.Dir("static/"))))

	log.Fatal(http.ListenAndServe(":8000", nil))
}
