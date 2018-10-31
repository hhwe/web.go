package main

import (
	"log"
	"net/http"
	"time"
)

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

func FindUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing finalHandler")
	time.Sleep(time.Second * 1)
	w.Write([]byte("OK"))
}

func FindUsersByID(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing finalHandler")
	vs := r.URL.Query()
	log.Println(r.URL.RawQuery)
	w.Write([]byte(vs.Get("id")))
}

func main() {
	app := NewApp(Pathing, Logging)
	app.AddRoute("/users", http.HandlerFunc(FindUsers))
	app.AddRoute("/users/:id", http.HandlerFunc(FindUsersByID))
	log.Fatal(http.ListenAndServe(":8080", app))
}
