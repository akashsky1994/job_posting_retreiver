package config

import (
	"fmt"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewConfig(configPaths ...string) *Config {
	v := viper.New()
	v.SetConfigName(getenv("JB_ENV", "development"))
	v.SetConfigType("yaml")
	v.AutomaticEnv()
	v.SetDefault("ServerPort", 8080)
	v.SetDefault("Env", "development")
	for _, path := range configPaths {
		v.AddConfigPath(path)
	}
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("failed to read the configuration file: %s", err))
	}
	var conf Config
	v.Unmarshal(&conf)
	return &conf
}

type Config struct {
	Router     *chi.Mux
	Logger     *log.Logger
	Cron       *cron.Cron
	Env        string
	ServerPort string
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
