package main

import (
	"net/http"
)

type Router struct {
	pattern        string
	params         map[int]string
	controllerType reflect.Type
}

func (r *Router) Match(r *http.Request) {

}

func (r *Router) findControllerInfo(r *http.Request) (string, string) {
	path := r.URL.Path
	if strings.HasSuffix(path, "/") {
		path = strings.TrimSuffix(path, "/")
	}
	pathInfo := strings.Split(path, "/")

	controllerName := defController
	if len(pathInfo) > 1 {
		controllerName = pathInfo[1]
	}

	methodName := defMethod
	if len(pathInfo) > 2 {
		methodName = strings.Title(strings.ToLower(pathInfo[2]))
	}

	return controllerName, methodName
}

func (r *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	controllerName, methodName := r.findControllerInfo(r)
	controllerT, ok := mapping[controllerName]
	if !ok {
		http.NotFound(w, r)
		return
	}

	refV := reflect.New(controllerT)
	method := refV.MethodByName(methodName)
	if !method.IsValid() {
		http.NotFound(w, r)
		return
	}

	controller := refV.Interface().(IApp)
	controller.Init(ctx)
	method.Call(nil)
}
