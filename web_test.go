package main

import (
	"fmt"
	"net/http"
	"testing"
)

type keyType int

const (
	key1 keyType = iota
	key2
)

func TestContext(t *testing.T) {
	assertEqual := func(val interface{}, exp interface{}) {
		if val != exp {
			t.Errorf("Expected %v, got %v.", exp, val)
		}
	}

	r, _ := http.NewRequest("GET", "http://localhost:8080/", nil)

	assertEqual(r.Response.StatusCode, 200)
}

func BenchmarkHello(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("hello")
	}
}
