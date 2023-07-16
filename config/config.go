package config

import (
	"fmt"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func NewConfig(env_type string, configPaths ...string) *Config {
	fmt.Println(env_type)
	v := viper.New()
	v.SetConfigName(getenv("JB_ENV", env_type))
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
	DB         *gorm.DB // the shared DB ORM object
	Router     *chi.Mux
	Logger     *log.Logger
	Cron       *cron.Cron
	Env        string
	ServerPort string
	DBPath     string
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
