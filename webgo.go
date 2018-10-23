package main

import (
	"net/http"
)

type App struct {
	handler     http.Handler
	middlewares []func(http.Handler) http.Handler
}

func Classic(h http.Handler) *App {
	return &App{
		handler:     h,
		middlewares: []func(http.Handler) http.Handler{middlewareTwo, middlewareOne},
	}
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, f := range app.middlewares {
		app.handler = f(app.handler)
	}
	app.handler.ServeHTTP(w, r)
}
