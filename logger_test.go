package webgo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogger(t *testing.T) {
	// Mock a app as http Handler
	app := NewApp(Logging)
	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "http://localhost:8080/foobar", nil)
	if err != nil {
		t.Error(err)
	}

	app.ServeHTTP(recorder, req)

	expect(t, recorder.Code, http.StatusNotFound)
}

func TestLogger_URLEncodedString(t *testing.T) {
	// Mock a app as http Handler
	app := NewApp(Logging)
	app.AddRoute("/users/:id", http.HandlerFunc((func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test logger url encoded string"))
	})))

	recorder := httptest.NewRecorder()

	// Test reserved characters - !*'();:@&=+$,/?%#[]
	req, err := http.NewRequest("GET", "http://localhost:8080/users/123232%0A%0A", nil)
	if err != nil {
		t.Error(err)
	}

	app.ServeHTTP(recorder, req)

	expect(t, recorder.Code, http.StatusOK)
}

func TestLogger_CustomFormat(t *testing.T) {
	// Mock a app as http Handler
	app := NewApp(Logging)
	app.AddRoute("/users", http.HandlerFunc((func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test logger url encoded string"))
	})))
	recorder := httptest.NewRecorder()

	userAgent := "Negroni-Test"
	req, err := http.NewRequest("GET", "http://localhost:8080/foobar?foo=bar", nil)
	if err != nil {
		t.Error(err)
	}
	req.Header.Set("User-Agent", userAgent)

	app.ServeHTTP(recorder, req)

	expect(t, recorder.Code, http.StatusNotFound)
}
