package main

import (
	"log"
	"net/http"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing middlewareOne")
		next.ServeHTTP(w, r)
		log.Println("Executing middlewareOne again")
	})
}

func Pathing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing middlewareTwo")
		// if r.URL.Path != "/" {
		// 	return
		// }
		next.ServeHTTP(w, r)
		log.Println("Executing middlewareTwo again")
	})
}

func final(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing finalHandler")
	w.Write([]byte("OK"))
}

func FindID(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing finalHandler")
	vs := r.URL.Query()
	log.Println("url query", vs)
	w.Write([]byte(vs.Get("id")))
}

func main() {
	app := NewApp(Logging, Pathing)
	app.AddRoute("/users", http.HandlerFunc(final))
	app.AddRoute("/users/:id", http.HandlerFunc(FindID))
	log.Fatal(http.ListenAndServe(":8080", app))
}
