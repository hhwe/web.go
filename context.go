package main

import (
	"net/http"
	"sync"
)



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
