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

func TestApp(t *testing.T) {
	// Mock a app as http Handler
	app := NewApp(Logging)
	app.AddRoute("/users", http.HandlerFunc(FindUsers))
	app.AddRoute("/users/:id", http.HandlerFunc(FindUsersByID))

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	app.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := "OK"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	req1, err := http.NewRequest("GET", "/users/123", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr1 := httptest.NewRecorder()
	app.ServeHTTP(rr1, req1)
	expected = "123"
	if rr1.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr1.Body.String(), expected)
	}
}

func TestApp_ServeHTTP(t *testing.T) {
	app := NewApp(Logging)

	req, err := http.NewRequest("GET", "/static", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	expect(t, rr.Code, http.StatusOK)
}
