package main

import (
	"net/http"
)

func main() {
	app := New()
	// if err != nil {
	// 	log.Fatalln("Could not start server:", err)
	// }
	app.AttachLogger()
	app.AttachRouter()
	app.AttachCron()
	app.FileServer()
	app.PrintRoutes()
	app.Logger.Info("Starting Job Retreiver App")
	app.Logger.Fatal(http.ListenAndServe(":"+app.ServerPort, app.Router))
}
