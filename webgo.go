// web.go is a micro restful framework with golang
// you can build a RESTful API with it.
// Also completed a middleware chain for a  request.
package main

import (
	"net/http"
)

type middleware func(http.Handler) http.Handler

// App is a HTTP multiplexer is a HTTP multiplexer / router similar to net/http.ServeMux.
type App struct {
	handler     http.Handler
	middlewares []middleware
	routes      map[string]*Router
	root        bool
}

// NewApp return a new app instance.
func NewApp() *App {
	return &App{}
}

// Classic return a basic app instance with some common middleware
func Classic(h http.Handler) *App {
	return &App{
		handler:     h,
		middlewares: []middleware{middlewareTwo, middlewareOne},
	}
}

// func (app *App) Handle(pattern string, handler http.Handler) {
// 	if pattern == "" {
// 		panic("http: invalid pattern")
// 	}
// 	if handler == nil {
// 		panic("http: nil handler")
// 	}

// 	if app.routes == nil {
// 		app.routes = []Router{}
// 	}

// 	if _, exist := app.routes[pattern]; exist {
// 		panic("http: multiple registrations for " + pattern)
// 	}

// 	route = make(map[string]Router)
// 	app.routes = append(app.routes, r)

// 	if pattern[0] != '/' {
// 		app.root = true
// 	}
// }

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, f := range app.middlewares {
		app.handler = f(app.handler)
	}
	app.handler.ServeHTTP(w, r)
}
