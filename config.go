package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

type AppConfig struct {
	Router     *chi.Mux
	Logger     *log.Logger
	Cron       *cron.Cron
	Env        string
	ServerPort string
}
