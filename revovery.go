package webgo

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				Logger.Error(fmt.Sprintf("provided ErrorHandlerFunc panic'd: %s, trace:\n%s", err, debug.Stack()))
				Logger.Error(fmt.Sprintf("%s\n", debug.Stack()))

				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
