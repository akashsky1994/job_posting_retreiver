package main

import (
	"flag"
	"net/http"
)

func main() {
	var env string
	flag.StringVar(&env, "env", "development", "sets environment type")
	flag.Parse()
	app := New(env)
	// if err != nil {
	// 	log.Fatalln("Could not start server:", err)
	// }

	app.AttachLogger()
	app.SetupDB()
	app.AttachRouter()
	app.AttachCron()
	app.FileServer()
	app.PrintRoutes()
	app.Logger.Info("Starting Job Retreiver App")
	app.Logger.Fatal(http.ListenAndServe(":"+app.ServerPort, app.Router))
}
