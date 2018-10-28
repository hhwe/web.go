package main

import (
	"net/http"
	"regexp"
)

// Route stores infomation about a matched route.
type Route struct {
	regex    *regexp.Regexp
	params   map[int]string
	handler http.Handler
}
