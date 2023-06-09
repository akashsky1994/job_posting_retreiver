package main

import (
	"akashsky1994/job_retreiver/handler"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

func New() (*AppConfig, error) {
	var err error
	var conf AppConfig

	conf.AttachLogger()
	conf.AttachRouter()
	conf.AttachCron()
	conf.FileServer()
	conf.PrintRoutes()
	return &conf, err
}

func (app *AppConfig) AttachLogger() error {
	logFile, err := os.OpenFile("job_retreiver.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
		defer logFile.Close()
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	app.Logger = log.New()
	app.Logger.SetOutput(mw)
	app.Logger.SetFormatter(&log.TextFormatter{ForceColors: true, FullTimestamp: true, TimestampFormat: "2006-01-02 15:04:05"})
	app.Logger.SetLevel(log.InfoLevel)
	app.Logger.Print("Logging to a job_retreiver.log in Go!")
	return nil
}

func (app *AppConfig) AttachRouter() {
	app.Router = chi.NewRouter()
	app.Router.Use(
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
		r.Mount("/v1", BuiltInRoutes())
	})
}

func (app *AppConfig) AttachCron() {
	app.Cron = cron.New()
	app.Cron.AddFunc("@daily", func() {
		app.Logger.Info("Added Cron")
		handler.FetchBuiltInJobs("147")
		handler.FetchBuiltInJobs("149")
	})
	app.Cron.Start()
	app.Logger.Infof("Cron Info: %+v\n", app.Cron.Entries())
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func (app *AppConfig) FileServer() {
	var path string = "/files"
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "data"))

	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	app.Router.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(filesDir))
		fs.ServeHTTP(w, r)
	})
}

func (app *AppConfig) PrintRoutes() {
	// Using Closures
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		app.Logger.Printf("%s %s\n", method, route) // Walk and print out all routes
		return nil
	}
	if err := chi.Walk(app.Router, walkFunc); err != nil {
		app.Logger.Panicf("Logging err: %s\n", err.Error()) // panic if there is an error
	}
}

func RespondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	fmt.Println(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
