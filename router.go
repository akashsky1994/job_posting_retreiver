package main

import (
	"job_posting_retreiver/handler"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/go-chi/chi/v5/middleware"
)

func (app *AppConfig) AttachRouter() {
	app.Router = chi.NewRouter()
	app.Router.Use(
		middleware.Heartbeat("/ping"),
		middleware.Throttle(1),
		middleware.RequestID,
		middleware.Logger,
		middleware.RedirectSlashes, // Redirect slashes to no slash URL versions
		middleware.Recoverer,       // Recover from panics without crashing server
		middleware.CleanPath,
	)
	app.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	app.Router.Route("/builtin", func(r chi.Router) {
		builtinhandler := handler.NewBuiltInHandler(app.Config)
		r.Mount("/v1", BuiltInRoutes(builtinhandler))
	})
}

func BuiltInRoutes(bihandler *handler.BuiltInHandler) *chi.Mux {
	router := chi.NewRouter()
	router.Group(func(router chi.Router) {
		router.Get("/jobs/{category_id}", bihandler.FetchJobsHandler)
	})
	return router
}
