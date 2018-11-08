// web.go is a micro restful framework with golang
// you can build a RESTful API with it.
// Also completed a middleware chain for a  request.
package webgo

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type Middleware func(http.Handler) http.Handler

// App is a HTTP multiplexer / router similar to net/http.Serveapp.
type App struct {
	middlewares []Middleware
	routes      map[string]Route
}

func NewApp(middlewares ...Middleware) *App {
	return &App{
		middlewares: middlewares,
		routes:      make(map[string]Route),
	}
}

// Route stores information to match a request and build URLs.
type Route struct {
	regexp *regexp.Regexp
	// todo: use reflect type mapping handler function
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
	Logger.Info(fmt.Sprintf("route mapping %s -> %s", pattern, handler))
}

// StaticWeb handle files from given file system root.
func (app *App) StaticWeb(pattern, path string) {
	staticHandler := http.StripPrefix(pattern, http.FileServer(http.Dir(path)))
	app.AddRoute(pattern, staticHandler)
}

// ServeHTTP dispatches the request to the handler whose
// pattern most closely matches the request URL.
func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// determines if the given path needs drop "/" to it.
	path := r.URL.Path
	if strings.HasPrefix(path, "/static") {
		http.StripPrefix("/static", http.FileServer(http.Dir("."))).ServeHTTP(w, r)
		return
	}
	if path != "/" && strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}

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

		// implement http serve with app's middlewares
		app.Handler(route.handler).ServeHTTP(w, r)
		return
	}

	// if no route matched, response with 404
	app.Handler(http.HandlerFunc(NotFound)).ServeHTTP(w, r)
}

// Run attaches the router to a http.Server and starts listening and serving HTTP requests.
// It is a shortcut for http.ListenAndServe(addr, router)
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (app *App) Run(addr ...string) {
	address := resolveAddress(addr)
	Logger.Info(fmt.Sprintf(`webgo restful api application 
* Environment: production
* Debug mode: off
* Running on http://%s/ (Press CTRL+C to quit)
`, address))
	Logger.Info(http.ListenAndServe(address, app))
}

func (app *App) Handler(h http.Handler) http.Handler {
	for _, m := range app.middlewares {
		h = m(h)
	}
	return h
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "404 page not found", http.StatusNotFound)
}
