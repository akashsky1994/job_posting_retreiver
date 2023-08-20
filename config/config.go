package config

import (
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/patrickmn/go-cache"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func NewConfig(env_file string, configPaths ...string) (*Config, error) {
	v := viper.New()
	v.SetConfigName(env_file)
	v.SetConfigType("env")

	v.SetDefault("SERVER_PORT", 8080)
	v.SetDefault("JB_ENV", "development")
	for _, path := range configPaths {
		v.AddConfigPath(path)
	}
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("failed to read the configuration file: %s", err))
	}
	var conf Config
	err := v.Unmarshal(&conf)
	if err != nil {
		return &conf, err
	}
	return &conf, nil
}

type Config struct {
	DB          *gorm.DB // the shared DB ORM object
	Router      *chi.Mux
	Logger      *log.Logger
	Cron        *cron.Cron
	Cache       *cache.Cache
	JB_ENV      string
	SERVER_PORT string
	DB_HOST     string
	DB_NAME     string
	DB_USER     string
	DB_PASSWORD string
	DB_PORT     string
	SENTRY_DSN  string
}

// func getenv(key, fallback string) string {
// 	value := os.Getenv(key)
// 	if len(value) == 0 {
// 		return fallback
// 	}
// 	return value
// }
