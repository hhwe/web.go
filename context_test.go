package webgo

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	// Mock a app as http Handler
	app := NewApp(Context)
	app.AddRoute("/status", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Get data from context
		if username := r.Context().Value("Username"); username != nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Hello " + username.(string) + "\n"))
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Logged in"))
		}
	}))

	// test request with no cookies
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost:8080/status", nil)
	if err != nil {
		t.Error(err)
	}
	app.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusNotFound)
	expect(t, recorder.Body.String(), "Not Logged in")

	// request with cookies
	recorder1 := httptest.NewRecorder()
	expiration := time.Now().Add(365 * 24 * time.Hour)
	req.AddCookie(&http.Cookie{Name: "username", Value: "alice_cooper@gmail.com", Expires: expiration})
	app.ServeHTTP(recorder1, req)
	expect(t, recorder1.Code, http.StatusOK)
	expect(t, recorder1.Body.String(), "Hello alice_cooper@gmail.com\n")
}
