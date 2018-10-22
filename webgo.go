package main

import "net/http"

type Handler interface {
	ServerHttp(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type middleware struct {
	handler Handler
	next *middleware
}

func (m middleware) ServerHttp(w http.ResponseWriter, r *http.Request) {
	m.handler.ServerHttp(w, r, m.next.ServerHttp)
}

type App struct {
	middleware middleware
	handlers []Handler
}

func (a *App) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	a.middleware.ServerHttp(w, r)
}

func NewApp(handlers ...Handler) *App {
	return &App{
		handlers:handlers,
		middleware:middleware{},
	}
}

func main() {
	mux := http.NewServeMux()
	n := NewApp(NewLogger())
	lh :=
	mux.Handle("/", )
}
