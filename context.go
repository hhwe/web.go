// context is based on go standard context package
package webgo

import (
	"context"
	"log"
	"net/http"
)

// Authentication use r.Context() make request context
func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, "-", r.RequestURI)
		cookie, _ := r.Cookie("username")
		if cookie != nil {
			//Add data to context
			ctx := context.WithValue(r.Context(), "Username", cookie.Value)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

type ResponseWriter struct {
	http.ResponseWriter
	status int
	// size   int
}

func (r *ResponseWriter) SetStatus(status int) {
	r.status = status
}

func (r *ResponseWriter) GetStatus() int {
	return r.status
}

// todo: pass context when handler every requests.
// Context is the most important part of webgo. It allows us to pass variables between middleware,
// manage the flow, validate the JSON of a request and render a JSON response for example.
type Context struct {
	Response    ResponseWriter
	Request     *http.Request
	Params      map[string]string
	Application *App
}

// WriteString writes string data into the response object.
func (ctx *Context) WriteString(content string) {
	ctx.Response.Write([]byte(content))
}

// Abort is a helper method that sends an HTTP header and an optional
// body. It is useful for returning 4xx or 5xx errors.
// Once it has been called, any return value from the handler will
// not be written to the response.
func (ctx *Context) Abort(status int, body string) {
	ctx.SetHeader("Content-Type", "application/json; charset=utf-8", true)
	ctx.Response.WriteHeader(status)
	ctx.Response.Write([]byte(body))
}

// Redirect is a helper method for 3xx redirects.
func (ctx *Context) Redirect(status int, url_ string) {
	ctx.Response.Header().Set("Location", url_)
	ctx.Response.WriteHeader(status)
	ctx.Response.Write([]byte("Redirecting to: " + url_))
}

//BadRequest writes a 400 HTTP response
func (ctx *Context) BadRequest() {
	ctx.Response.WriteHeader(400)
}

// Notmodified writes a 304 HTTP response
func (ctx *Context) NotModified() {
	ctx.Response.WriteHeader(304)
}

//Unauthorized writes a 401 HTTP response
func (ctx *Context) Unauthorized() {
	ctx.Response.WriteHeader(401)
}

//Forbidden writes a 403 HTTP response
func (ctx *Context) Forbidden() {
	ctx.Response.WriteHeader(403)
}

// NotFound writes a 404 HTTP response
func (ctx *Context) NotFound(message string) {
	ctx.Response.WriteHeader(404)
	ctx.Response.Write([]byte(message))
}

// SetHeader sets a response header. If `unique` is true, the current value
// of that header will be overwritten . If false, it will be appended.
func (ctx *Context) SetHeader(hdr string, val string, unique bool) {
	if unique {
		ctx.Response.Header().Set(hdr, val)
	} else {
		ctx.Response.Header().Add(hdr, val)
	}
}

// SetCookie adds a cookie header to the response.
func (ctx *Context) SetCookie(cookie *http.Cookie) {
	ctx.SetHeader("Set-Cookie", cookie.String(), false)
}
