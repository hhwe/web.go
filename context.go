package main

import (
	"net/http"
	"sync"
	"time"
)

// Context --------------------------------------------------------------------

// app context split by request
type Context struct {
	Data map[*http.Request]map[interface{}]interface{}
	sync.RWMutex
}

var context = Context{Data: make(map[*http.Request]map[interface{}]interface{})}

// Set stores a value for a given key in a given request.
func (c *Context) Set(r *http.Request, key, val interface{}) {
	c.Lock()
	if c.Data[r] == nil {
		c.Data[r] = make(map[interface{}]interface{})
	}
	c.Data[r][key] = val
	c.Unlock()
}

// Get returns a value stored for a given key in a given request.
func (c *Context) Get(r *http.Request, key interface{}) interface{} {
	c.RLock()
	if ctx := c.Data[r]; ctx != nil {
		value := ctx[key]
		c.RUnlock()
		return value
	}
	c.RUnlock()
	return nil
}

// Delete removes a value stored for a given key in a given request.
func (c *Context) Delete(r *http.Request, key interface{}) {
	c.Lock()
	if c.Data[r] != nil {
		delete(c.Data[r], key)
	}
	c.Unlock()
}

// Clear removes all values stored for a given request.
func (c *Context) Clear(r *http.Request) {
	c.Lock()
	c.clear(r)
	c.Unlock()
}

// clear is Clear without the lock.
func (c *Context) clear(r *http.Request) {
	delete(c.Data, r)
}


// Session --------------------------------------------------------------------

//// Session stores the values and optional configuration for a session.
//type Session struct {
//	ID     string                      // The ID of the session.
//	Values map[interface{}]interface{} // Values contains in the session.
//}
//
//// GetRegistry returns a registry instance for the current request.
//func GetRegistry(r *http.Request) *Registry {
//	registry := context.Get(r, "session")
//	if registry != nil {
//		return registry.(*Registry)
//	}
//	newRegistry := &Registry{
//		request:  r,
//		sessions: make(map[string]sessionInfo),
//	}
//	context.Set(r, registryKey, newRegistry)
//	return newRegistry
//}

var sessions = map[string]string{}


// Cookie ---------------------------------------------------------------------

// addCookie will apply a new cookie to the response of a http
// request, with the key/value this method is passed.
func addCookie(w http.ResponseWriter, name string, value string) {
	expire := time.Now().AddDate(0, 0, 3)
	cookie := http.Cookie{
		Name:    name,
		Value:   value,
		Expires: expire,
		Path:    "/",
	}
	http.SetCookie(w, &cookie)
}

// validate cookie from request of current request
func validCookie(w http.ResponseWriter, r *http.Request, name string) {
	cookie, err := r.Cookie(name)
	if err != nil {
		panic(err)
	}

	if cookie.Value == "" || sessions[cookie.Value] != "" {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	}
}
