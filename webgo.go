package main

import (
	"net/http"
	"reflect"
)

type middleware func(http.Handler) http.Handler

// App is a HTTP multiplexer is a HTTP multiplexer / router similar to net/http.ServeMux.
type App struct {
	handler     http.Handler
	middlewares []middleware
	routes      map[string]Router
	root        bool
}


func Classic(h http.Handler) *App {
	return &App{
		handler:     h,
		middlewares: []func(http.Handler) http.Handler{middlewareTwo, middlewareOne},
	}
}

func (app *App) Handle(pattern string, handler http.Handler) {
	if pattern == "" {
		panic("http: invalid pattern")
	}
	if handler == nil {
		panic("http: nil handler")
	}

	if app.routes == nil {
		app.routes = []Router
	}

	if _, exist := app.routes[pattern]; exist {
		panic("http: multiple registrations for " + pattern)
	}

	route = make(map[string]Router)
	app.routes = append(app.routes, r)

	if pattern[0] != '/' {
		app.root = true
	}
}


func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, f := range app.middlewares {
		app.handler = f(app.handler)
	}
	app.handler.ServeHTTP(w, r)
}
