package main

import (
	"log"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"
)

// Route stores infomation about a matched route.
type Route struct {
	regex    *regexp.Regexp
	params   map[int]string
	resource reflect.Type
}

type Router struct {
	routers     []*Route
	Application *App
}

// Add is mapping dynamic url to resource
func (rt *Router) Add(pattern string, rc Resource) {
	parts := strings.Split(pattern, "/")
	j, params := 0, make(map[int]string)
	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			// replace params with regexp
			expr := "([^/]+)"
			parts[i] = expr
			params[j] = part[1:]
			j++
		}
	}

	// now create the Route
	// recreate the url pattern, with parameters replaced
	// by regular expressions. then compile the regex
	pattern = strings.Join(parts, "/")
	regex := regexp.MustCompile(pattern)
	t := reflect.Indirect(reflect.ValueOf(rc)).Type()
	route := &Route{
		regex:    regex,
		params:   params,
		resource: t,
	}
	rt.routers = append(rt.routers, route)
}

// ServeHTTP implement http.Handler
func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var started bool
	requestPath := r.URL.Path

	//find a matching Route
	for _, route := range rt.routers {

		//check if Route pattern matches url
		if !route.regex.MatchString(requestPath) {
			continue
		}

		//get submatches (params)
		matches := route.regex.FindStringSubmatch(requestPath)

		//double check that the Route matches the URL pattern.
		if len(matches[0]) != len(requestPath) {
			continue
		}

		params := make(map[string]string)
		if len(route.params) > 0 {
			//add url parameters to the query param map
			values := r.URL.Query()
			for i, match := range matches[1:] {
				values.Add(route.params[i], match)
				params[route.params[i]] = match
			}

			r.URL.RawQuery = url.Values(values).Encode() + "&" + r.URL.RawQuery
		}
		//Invoke the request handler
		vc := reflect.New(route.resource)
		init := vc.MethodByName("Init")
		in := make([]reflect.Value, 2)
		ct := &Context{ResponseWriter: w, Request: r, Params: params}
		in[0] = reflect.ValueOf(ct)
		in[1] = reflect.ValueOf(route.resource.Name())
		init.Call(in)
		in = make([]reflect.Value, 0)

		method := vc.MethodByName(strings.Title(r.Method))
		method.Call(in)

		break
	}

	//if no matches to url, throw a not found exception
	if started == false {
		http.NotFound(w, r)
	}
}

func middlewareRouter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing middleware router")
		if r.URL.Path != "/" {
			return
		}
		next.ServeHTTP(w, r)
		log.Println("Executed middleware router")
	})
}
