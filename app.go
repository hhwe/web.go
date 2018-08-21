package main

import (
	"fmt"
	"log"

	"net/http"
)

func main() {
	http.HandleFunc("/", Index)
	http.HandleFunc("/movies1", AllMovies)
	http.HandleFunc("/movies2", CreateMovie)
	http.HandleFunc("/movies3", UpdateMovie)
	http.HandleFunc("/movies4", DeleteMovie)
	http.HandleFunc("/movies/{id}", FindMovie)

	log.Fatal(http.ListenAndServe(":8888", nil))
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World!")
}

func AllMovies(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented yet !")
}

func FindMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented yet !")
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented yet !")
}

func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented yet !")
}

func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "not implemented yet !")
}
