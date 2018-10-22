package main

import (
	"log"
	"net/http"
	"time"
)

type timeHandler struct {
	format string
}

func (th *timeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tm := time.Now().Format(th.format)
	w.Write([]byte("The time is:" + tm))
}

func timeHandler1(w http.ResponseWriter, r *http.Request) {
	tm := time.Now().Format(time.RFC1123)
	w.Write([]byte("secodnigjdfk The time is: " + tm))
}

func main() {
	mux := http.NewServeMux()

	rh := http.RedirectHandler("http://example.org", 307)
	mux.Handle("/foo", rh)

	th := &timeHandler{format: time.RFC1123}
	mux.Handle("/time", th)

	mh := http.TimeoutHandler(th, time.Nanosecond*1, "timeout")
	mux.Handle("/out", mh)

	th1 := http.HandlerFunc(timeHandler1)
	mux.Handle("/time1", th1)

	log.Println("Listening...")
	http.ListenAndServe(":3000", mux)
}
