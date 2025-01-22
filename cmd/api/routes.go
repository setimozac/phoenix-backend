package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	// middlewares
	mux.Use(middleware.Recoverer)
	mux.Use(app.enableCORS)

	mux.Get("/", app.Hello)
	mux.Get("/services", app.GetAllEnvManagers)
	// mux.Get("/test-add-get", app.TestGetEnvManager)

	return mux
}