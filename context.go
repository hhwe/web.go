package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"strconv"
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

// todo: store session on redis, add expire time to every sessions
const (
	sessionToken = iota
	sessionCookie
	sessionCRSF
)
var sessions = map[int]map[string]string{}


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
		HttpOnly:true,  // can't called by javascript
	}
	http.SetCookie(w, &cookie)
	sessions[sessionCookie] = map[string]string{name:value}
}

// validate cookie from request of current request
func validCookie(r *http.Request, name string) bool {
	cookie, err := r.Cookie(name)
	if err != nil {
		panic(err)
	}

	// check cookie's value is valid and in sessions
	if sessions[sessionCookie][cookie.Name] == cookie.Value {
		logger.Println("Warning: login error")
		return true
	}
	return false
}


// Token ----------------------------------------------------------------------

// Generates a random number for prevention CSRF
func addToken(name string) (token string) {
	h := md5.New()
	currentTime := time.Now().Unix()
	io.WriteString(h, strconv.FormatInt(currentTime, 10))
	io.WriteString(h, HashSha256("token"))
	token = fmt.Sprintf("%x", h.Sum(nil))
	sessions[sessionToken] =map[string]string{name:token}
	return
}
