// web.go is a micro restful framework with golang
// you can build a RESTful API with it.
// Also completed a middleware chain for a  request.
package main

import (
	"net/http"
)

type middleware func(http.Handler) http.Handler

// App is a HTTP multiplexer / router similar to net/http.ServeMux.
type App struct {
	middlewares []middleware
	routes      []*Route
}

// NewApp registers an empty route.
func (a *App) NewApp() *Route {
	route := &Route{}
	a.routes = append(a.routes, route)
	return route
}

// Handle registers a new route with a matcher for the URL path.
func (a *App) Handle(path string, handler http.Handler) {
	return a.NewApp().handler
}

func (app *App) Match(r *http.Request) bool {
	for _, route := range app.routes {
		if route.Match(r, m) {

		}
	}
	return true
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//app.handler = app.routes
	var handler http.Handler
	var route Route
	if app.Match(r) {
		handler = route.resource
	}

	for _, f := range app.middlewares {
		handler = f(handler)
	}

	handlea.ServeHTTP(w, r)
}
