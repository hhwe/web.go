package main

import (
	"net/http"
	"sync"
	"time"
)

// app context split by request
type Context struct {
	Data map[*http.Request]map[interface{}]interface{}
	sync.RWMutex
}

var context = Context{Data:make(map[*http.Request]map[interface{}]interface{})}

// Set stores a value for a given key in a given request.
func (c *Context)Set(r *http.Request, key, val interface{}) {
	c.Lock()
	if c.Data[r] == nil {
		c.Data[r] = make(map[interface{}]interface{})
	}
	c.Data[r][key] = val
	c.Unlock()
}

// Get returns a value stored for a given key in a given request.
func (c *Context)Get(r *http.Request, key interface{}) interface{} {
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
func (c *Context)Delete(r *http.Request, key interface{}) {
	c.Lock()
	if c.Data[r] != nil {
		delete(c.Data[r], key)
	}
	c.Unlock()
}

// Clear removes all values stored for a given request.
func (c *Context)Clear(r *http.Request) {
	c.Lock()
	c.clear(r)
	c.Unlock()
}

// clear is Clear without the lock.
func (c *Context)clear(r *http.Request) {
	delete(c.Data, r)
}

// session stored in redis
type Session interface {
	Set(key, value interface{}) error // set session value
	Get(key interface{}) interface{}  // get session value
	Delete(key interface{}) error     // delete session value
	SessionID() string                // back current sessionID
}

// addCookie will apply a new cookie to the response of a http
// request, with the key/value this method is passed.
func addCookie(w http.ResponseWriter, name string, value string) {
	expire := time.Now().AddDate(0, 0, 3)
	cookie := http.Cookie{
		Name:    name,
		Value:   value,
		Expires: expire,
	}
	http.SetCookie(w, &cookie)
}

// validate cookie from request of current user
func validCookie(r *http.Request, name string, value string) bool {
	cookie, err := r.Cookie(name)
	if err != nil {
		panic(err)
	}
	if cookie.Value == value {
		return true
	}
	return false
}