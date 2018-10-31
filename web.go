// web.go is a micro restful framework with golang
// you can build a RESTful API with it.
// Also completed a middleware chain for a  request.
package main

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type middleware func(http.Handler) http.Handler

// App is a HTTP multiplexer / router similar to net/http.Serveapp.
type App struct {
	middlewares []middleware
	routes      map[string]Route
}

func NewApp(middlewares ...middleware) *App {
	return &App{
		middlewares: middlewares,
		routes:      make(map[string]Route),
	}
}

// Route stores information to match a request and build URLs.
type Route struct {
	regexp  *regexp.Regexp
	handler http.Handler
	params  []string
}

// AddRoute registers the handler for the given pattern.
func (app *App) AddRoute(pattern string, handler http.Handler) {
	if _, exist := app.routes[pattern]; exist {
		panic("http: multiple registrations for " + pattern)
	}

	parts := strings.Split(pattern, "/")
	re := "([^/]+)"
	params := make([]string, 0)
	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			parts[i] = re
			params = append(params, part[1:len(part)])
		}
	}

	if app.routes == nil {
		app.routes = make(map[string]Route)
	}
	app.routes[pattern] = Route{
		regexp:  regexp.MustCompile(strings.Join(parts, "/")),
		handler: handler,
		params:  params,
	}
}

// ServeHTTP dispatches the request to the handler whose
// pattern most closely matches the request URL.
func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// find a matching route
	path := r.URL.Path
	for _, route := range app.routes {
		matches := route.regexp.FindStringSubmatch(path)
		if matches == nil || len(matches[0]) != len(path) {
			continue
		}
		if len(route.params) != len(matches)-1 {
			continue
		}

		// add parameter to url encoding
		values := r.URL.Query()
		for i, match := range matches[1:] {
			values.Add(route.params[i], match)
		}
		r.URL.RawQuery = url.Values(values).Encode() + "&" + r.URL.RawQuery

		// appliment http serve
		app.hanlder(route.handler).ServeHTTP(w, r)
		break
	}
}

func (app *App) hanlder(h http.Handler) http.Handler {
	for _, m := range app.middlewares {
		h = m(h)
	}
	return h
}
