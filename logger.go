package webgo

import (
	"log"
	"net/http"
	"os"
	"time"

	"webgo/logging"
)

var Logger = logging.NewLogger(os.Stderr, log.LstdFlags)

func init() {
	Logger.SetLever("debug")
	Logger.SetColor(true)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		Logger.Info(time.Since(start))
	})
}
