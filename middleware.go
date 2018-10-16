// Using closures to implement the middleware.
package main

import (
	"net/http"
	"time"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

// Logging logs all requests with its path and the time it took to process
func Logging() Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			// Do middleware things
			start := time.Now()
			defer func() {
				logger.Println(r.Method, r.URL.Path, time.Since(start))
			}()

			defer func() {
				e := recover()
				if e != nil {
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				}
			}()

			// Call the next middleware/handler in chain
			f(w, r)

		}
	}
}

// DBSession used to add a mongodb session in request context,
// make sure a request process have a isolated database session
func DBSession() Middleware {
	// return the Adapter
	return func(f http.HandlerFunc) http.HandlerFunc {
		// the adapter (when called) should return a new handler
		return func(w http.ResponseWriter, r *http.Request) {
			// copy the database session
			session := globalSession.Copy()
			defer session.Close() // clean up
			db := session.DB(database)

			// save it in the mux context, add to request lifetime
			context.Set(r, "database", db)

			// pass execution to the original handler
			f(w, r)

			// clears request values at the end of a request lifetime.
			context.Clear(r)
		}
	}
}

// Method ensures that url can only be requested with a specific method,
// else returns a 400 Bad Request
func Method(ms ...string) Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			// Do middleware things
			method := false
			for _, m := range ms {
				if r.Method == m {
					method = true
					break
				}
			}
			if !method {
				http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
				return
			}

			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}

// Check request's authorizations, make sure all important data save
func Auth() Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {
			if !validCookie(r, "name") {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

			// Call the next middleware/handler in chain
			f(w, r)

		}
	}
}

// Chain applies middlewares to a http.HandlerFunc
func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}
