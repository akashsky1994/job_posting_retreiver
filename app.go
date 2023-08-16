package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"job_posting_retreiver/config"
	"job_posting_retreiver/constant"
	"job_posting_retreiver/errors"
	"job_posting_retreiver/handler"
	"job_posting_retreiver/model"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/getsentry/raven-go"
	"github.com/go-chi/chi/v5"
	"github.com/patrickmn/go-cache"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AppConfig struct {
	*config.Config
}

func New(env string) (*AppConfig, error) {
	conf, err := config.NewConfig(env, "./config")
	if err != nil {
		return nil, err
	}
	return &AppConfig{
		conf,
	}, nil
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
		for _, source := range constant.JOB_SOURCES {
			app.Logger.Info("Running Job Retreiver Cron for:", source)
			jobhandler := handler.NewJobSourceHandler("trueup", app.Config)
			if err := jobhandler.FetchJobs(); err != nil {
				raven.CaptureErrorAndWait(err, map[string]string{"type": "cron_" + source})
				_, severity := errors.GetTypeAndLogLevel(err)
				app.Logger.Log(severity, err)

			}
			if err := jobhandler.ProcessJobs(); err != nil {
				raven.CaptureErrorAndWait(err, map[string]string{"type": "cron_" + source})
				_, severity := errors.GetTypeAndLogLevel(err)
				app.Logger.Log(severity, err)
			}
		}
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
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=America/New_York", app.DB_HOST, app.DB_USER, app.DB_PASSWORD, app.DB_NAME, app.DB_PORT)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{CreateBatchSize: 1})

	db = db.Session(&gorm.Session{CreateBatchSize: 1})
	if err != nil {
		app.Logger.Panicln("failed to connect database", err.Error())
	}
	app.DB = db

	// Migration
	err = app.Config.DB.AutoMigrate(
		&model.JobListing{},
		&model.Company{},
		&model.FileLog{},
		&model.Country{},
		&model.State{},
		&model.City{},
	)
	if err != nil {
		app.Logger.Error("Error Runing Automigration", err.Error())
	}
	err = app.LoadRegions()
	if err != nil {
		app.Logger.Error("Error Loading Region Data", err.Error())
	}
}

func (app *AppConfig) StartCache() {
	app.Cache = cache.New(5*time.Hour, 10*time.Hour)
}

func (app *AppConfig) LoadRegions() error {
	// Let's first read the `config.json` file
	content, err := ioutil.ReadFile("./data/countries+states+cities.json")
	if err != nil {
		return errors.DataProcessingError.Wrap(err, "Error Adding Region Data to DB | Error when opening file", log.ErrorLevel)
	}
	// Now let's unmarshall the data into `payload`
	var countries []model.Country
	err = json.Unmarshal(content, &countries)
	if err != nil {
		return errors.DataProcessingError.Wrap(err, "Error Adding Region Data to DB | Error during Unmarshal():", log.ErrorLevel)
	}

	// err = app.Config.DB.Create(&countries).Error
	// if err != nil {
	// 	return errors.DataProcessingError.Wrap(err, "Error Adding Region Data to DB", log.ErrorLevel)
	// }
	var AvailableRegions []string
	for _, country := range countries {
		if country.Name == "United States" {
			AvailableRegions = append(AvailableRegions, country.Name, country.ISO2, country.ISO3, country.Region)
			for _, state := range country.States {
				AvailableRegions = append(AvailableRegions, state.Name, state.StateCode)
				for _, city := range state.Cities {
					AvailableRegions = append(AvailableRegions, city.Name)
				}
			}
		}
	}
	app.Cache.Set("allowed_regions", AvailableRegions, cache.NoExpiration)
	return nil
}
