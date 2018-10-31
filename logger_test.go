package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

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
	app.AddRoute("/users", http.HandlerFunc(FindUsers))
	app.AddRoute("/users/:id", http.HandlerFunc(FindUsersByID))

	recorder := httptest.NewRecorder()

	// Test reserved characters - !*'();:@&=+$,/?%#[]
	req, err := http.NewRequest("GET", "http://localhost:8080/users%2F123232%0A%0A", nil)
	if err != nil {
		t.Error(err)
	}

	app.ServeHTTP(recorder, req)

	expect(t, recorder.Code, http.StatusOK)
}

func TestLogger_CustomFormat(t *testing.T) {
	// Mock a app as http Handler
	app := NewApp(Logging)
	app.AddRoute("/users", http.HandlerFunc(FindUsers))
	app.AddRoute("/users/:id", http.HandlerFunc(FindUsersByID))

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
