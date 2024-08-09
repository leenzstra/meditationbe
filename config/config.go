package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	PgUrl     string `mapstructure:"pg_url"`
}

var config *Config

func Init(env string) {
	var err error

	cfg := viper.New()
	cfg.SetConfigType("yaml")
	cfg.AddConfigPath("config/")
	cfg.SetConfigName(env)

	err = cfg.ReadInConfig()
	if err != nil {
		log.Fatal("error on parsing env configuration file")
	}

	err = cfg.Unmarshal(config)
	if err != nil {
		log.Fatal("error on parsing env configuration file")
	}
}


func GetConfig() *Config {
	return config
}