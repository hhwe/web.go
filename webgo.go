package main

import (
	"github.com/hhwe/webgo/handles"
)

type Webgo struct {
	http.Server
	handles []Handler
}

func (s *Webgo) ServerHTTP(w http.ResponseWriter, r *http.Request) {

}
