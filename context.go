package main

import "net/http"

type HandlersChain int

type Context struct {
	Response http.ResponseWriter
	Request *http.Request
	Handlers HandlersChain
}