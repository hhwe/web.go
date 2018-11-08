// Represents an abstract RESTful resource. Concrete resources should
// extend from this class and expose methods for each supported HTTP
// method. If a resource is invoked with an unsupported HTTP method,
// the API will return a response with status 405 Method Not Allowed.
// Otherwise the appropriate method is called and passed all arguments
// from the url rule used when adding the resource to an Api instance.
package webgo

import (
	"net/http"
	"reflect"
)

type Resourcer interface {
	DispatchRequest(w http.ResponseWriter, r *http.Request)
}

// AddResource registers the resource for the given pattern.
func (app *App) AddResource(pattern string, rs Resourcer) {
	rsp := reflect.ValueOf(rs)
	tp := rsp.Type()
	getMethod, ok := tp.MethodByName("Get")
	if ok {
		getHandler := 
		app.AddRoute(pattern, getHandler)
	}
}
