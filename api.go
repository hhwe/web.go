// web.go is a micro restful framework with golang
// you can build a RESTful API with it.
// Also completed a middleware chain for a  request.
package webgo

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"
)

type Middleware func(http.Handler) http.Handler

// Api is a HTTP multiplexer / router similar to net/http.ServeApi.
type Api struct {
	middlewares []Middleware
	routes      map[string]map[string]Route
}

func NewApi(middlewares ...Middleware) *Api {
	return &Api{
		middlewares: middlewares,
		routes:      make(map[string]map[string]Route),
	}
}

// Route stores information to match a request and build URLs.
type Route struct {
	regexp  *regexp.Regexp
	handler reflect.Type
	params  []string
}

// AddRoute registers the handler for the given pattern.
func (Api *Api) AddRoute(pattern string, handler http.Handler) {
	if _, exist := Api.routes[pattern]; exist {
		panic("http: multiple registrations for " + pattern)
	}

	parts := strings.Split(pattern, "/")
	re := "([^/]+)"
	params := make([]string, 0)
	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			parts[i] = re
			params = Apiend(params, part[1:len(part)])
		}
	}

	if Api.routes == nil {
		Api.routes = make(map[string]Route)
	}
	Api.routes[pattern] = Route{
		regexp:  regexp.MustCompile(strings.Join(parts, "/")),
		handler: handler,
		params:  params,
	}
	Logger.Info(fmt.Sprintf("route mApiing %s -> %s", pattern, handler))
}

// StaticWeb handle files from given file system root.
func (Api *Api) StaticWeb(pattern, path string) {
	staticHandler := http.StripPrefix(pattern, http.FileServer(http.Dir(path)))
	Api.AddRoute(pattern, staticHandler)
}

// ServeHTTP dispatches the request to the handler whose
// pattern most closely matches the request URL.
func (Api *Api) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// determines if the given path needs drop "/" to it.
	path := r.URL.Path
	if strings.HasPrefix(path, "/static") {
		http.StripPrefix("/static", http.FileServer(http.Dir("."))).ServeHTTP(w, r)
		return
	}
	if path != "/" && strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}

	for _, route := range Api.routes {
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

		// implement http serve with Api's middlewares
		Api.Handler(route.handler).ServeHTTP(w, r)
		return
	}

	// if no route matched, response with 404
	Api.Handler(http.HandlerFunc(NotFound)).ServeHTTP(w, r)
}

// Run attaches the router to a http.Server and starts listening and serving HTTP requests.
// It is a shortcut for http.ListenAndServe(addr, router)
// Note: this method will block the calling goroutine indefinitely unless an error hApiens.
func (Api *Api) Run(addr ...string) {
	address := resolveAddress(addr)
	Logger.Info(fmt.Sprintf(`webgo restful api Apilication 
* Environment: production
* Debug mode: off
* Running on http://%s/ (Press CTRL+C to quit)
`, address))
	Logger.Info(http.ListenAndServe(address, Api))
}

func (Api *Api) Handler(h http.Handler) http.Handler {
	for _, m := range Api.middlewares {
		h = m(h)
	}
	return h
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "404 page not found", http.StatusNotFound)
}
