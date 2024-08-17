package config

import (
	"log"

	"github.com/spf13/viper"
)

type EnvType string

const (
	Development EnvType = "development"
	Production  EnvType = "production"
	Local       EnvType = "local"
)

type Config struct {
	PgUrl     string `mapstructure:"MED_PG_URL"`
	Env       string `mapstructure:"MED_ENV"`
	JWTSecret string `mapstructure:"MED_JWT_SECRET"`

	S3Endpoint  string `mapstructure:"MED_S3_ENDPOINT"`
	S3AccessKey string `mapstructure:"MED_S3_ACCESS_KEY"`
	S3SecretKey string `mapstructure:"MED_S3_SECRET_KEY"`
	S3BucketId  string `mapstructure:"MED_S3_BUCKET_ID"`
	S3JWT       string `mapstructure:"MED_S3_JWT"`
}

var config *Config

func Init(env EnvType) {
	var err error
	cfg := viper.New()

	cfg.SetConfigType("env")
	cfg.AddConfigPath("config/")
	cfg.SetConfigName(string(env))
	cfg.AutomaticEnv()

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
