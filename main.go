package main

import (
	"flag"
	"net/http"

	"github.com/sirupsen/logrus"
)

func main() {
	var env string
	flag.StringVar(&env, "env", ".env", "provides environment filename")
	flag.Parse()
	app, err := New(env)
	if err != nil {
		logrus.Fatalln("Could not start server:", err)
	}

	app.AttachLogger()
	app.StartCache()
	app.SetupDB()
	app.AttachRouter()
	app.AttachCron()
	app.FileServer()
	app.PrintRoutes()
	app.Logger.Info("Starting Job Retreiver App")
	app.Logger.Fatal(http.ListenAndServe(":"+app.SERVER_PORT, app.Router))
}
