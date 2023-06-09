package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func main() {
	app, err := New()
	if err != nil {
		log.Fatalln("Could not start server:", err)
	}
	app.Logger.Info("Starting Job Retreiver App")
	app.ServerPort = "8080"
	app.Logger.Fatal(http.ListenAndServe(":"+app.ServerPort, app.Router))
}
