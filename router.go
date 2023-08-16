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
		jobhandler := handler.NewJobSourceHandler("builtin", app.Config)
		r.Mount("/v1", BuiltInRoutes(jobhandler))
	})
	app.Router.Route("/simplify", func(r chi.Router) {
		jobhandler := handler.NewJobSourceHandler("simplify", app.Config)
		r.Mount("/v1", SimplifyRoutes(jobhandler))
	})
	app.Router.Route("/trueup", func(r chi.Router) {
		jobhandler := handler.NewJobSourceHandler("trueup", app.Config)
		r.Mount("/v1", TrueUpRoutes(jobhandler))
	})
	app.Router.Route("/jobs", func(r chi.Router) {
		jobhandler := handler.NewJobHandler(app.Config)
		r.Mount("/v1", JobRoutes(jobhandler))
	})
}

func BuiltInRoutes(jobhandler handler.JobSourceHandler) *chi.Mux {
	router := chi.NewRouter()
	router.Group(func(router chi.Router) {
		router.Get("/fetch", jobhandler.AggregateJobs)
	})
	return router
}

func SimplifyRoutes(jobhandler handler.JobSourceHandler) *chi.Mux {
	router := chi.NewRouter()
	router.Group(func(router chi.Router) {
		router.Get("/fetch", jobhandler.AggregateJobs)
	})
	return router
}

func TrueUpRoutes(jobhandler handler.JobSourceHandler) *chi.Mux {
	router := chi.NewRouter()
	router.Group(func(router chi.Router) {
		router.Get("/fetch", jobhandler.AggregateJobs)
	})
	return router
}

func JobRoutes(jobhandler *handler.JobHandler) *chi.Mux {
	router := chi.NewRouter()
	router.Group(func(router chi.Router) {
		router.Get("/list", jobhandler.ListJobs)
		router.Post("/add", jobhandler.AddJobs)
	})
	return router
}
