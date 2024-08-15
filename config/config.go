package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	PgUrl     string `mapstructure:"pg_url"`
	Env       string `mapstructure:"env"`
	JWTSecret string `mapstructure:"jwt_secret"`
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
		log.Fatalf("error on parsing env configuration file: %v", err)
	}

	err = cfg.Unmarshal(&config)
	if err != nil {
		log.Fatalf("error on unmarshaling env configuration file: %v", err)
	}

	log.Printf("%v", config)
}

func GetConfig() *Config {
	return config
}
