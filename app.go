package main

import (
	"io"
	"job_posting_retreiver/config"
	"job_posting_retreiver/handler"
	"job_posting_retreiver/model"
	"net/http"
	"os"
	"strings"

	"github.com/glebarez/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AppConfig struct {
	*config.Config
}

func New(env string) *AppConfig {
	return &AppConfig{
		config.NewConfig(env, "./config"),
	}
}

func (app *AppConfig) AttachLogger() error {
	logFile, err := os.OpenFile("job_posting_retreiver.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
		defer logFile.Close()
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	app.Logger = log.New()
	app.Logger.SetOutput(mw)
	app.Logger.SetFormatter(&log.TextFormatter{ForceColors: true, FullTimestamp: true, TimestampFormat: "2006-01-02 15:04:05"})
	app.Logger.SetLevel(log.InfoLevel)
	app.Logger.Print("Logging to a job_posting_retreiver.log in Go!")
	return nil
}

func (app *AppConfig) AttachCron() {
	app.Cron = cron.New()
	app.Cron.AddFunc("@daily", func() {
		app.Logger.Info("BuiltIn Job Retreiver Cron Added")
		builtinhandler := handler.NewBuiltInHandler(app.Config)
		builtinhandler.FetchJobs("147")
		builtinhandler.FetchJobs("149")

		// app.Logger.Info("Simplify Job Retreiver Cron Added")
		// simplifyhandler := handler.NewSimplifyHandler(app.Config)
		// simplifyhandler.FetchJobs()

		app.Logger.Info("Trueup Job Retreiver Cron Added")
		trueuphandler := handler.NewTrueupHandler(app.Config)
		trueuphandler.FetchJobs()

	})
	app.Cron.Start()
	app.Logger.Infof("Cron Info: %+v\n", app.Cron.Entries())
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func (app *AppConfig) FileServer() {
	filesDir := http.Dir("data")

	app.Router.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
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

func (app *AppConfig) SetupDB() {
	db, err := gorm.Open(sqlite.Open(app.DBPath), &gorm.Config{
		CreateBatchSize: 1,
	})
	db = db.Session(&gorm.Session{CreateBatchSize: 1})
	if err != nil {
		app.Logger.Panicf("failed to connect database", err.Error())
	}
	app.DB = db

	// Migration
	err = app.Config.DB.AutoMigrate(
		&model.JobListing{},
		&model.Company{},
		&model.FileLogs{},
	)
	if err != nil {
		app.Logger.Error("Error Runing Automigration", err.Error())
	}
}
