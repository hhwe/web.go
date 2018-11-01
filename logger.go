package webgo

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var Logger = NewLogger(os.Stdout, log.LstdFlags)

func NewLogger(out io.Writer, flag int) *log.Logger {
	return log.New(out, "[web.go]", flag)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		Logger.Println(time.Since(start), r.URL.Path)
	})
}
