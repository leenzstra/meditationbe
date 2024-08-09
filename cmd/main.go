package main

import (
	"github.com/leenzstra/meditationbe/config"
	"github.com/leenzstra/meditationbe/db"
	"github.com/leenzstra/meditationbe/server"
	"go.uber.org/zap"
	"log"
)

func Run() {
	var err error

	config.Init("development")
  cfg := config.GetConfig()

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}

	db, err := db.NewPostgres(cfg.PgUrl, logger)
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	server.NewRouter(db, logger).Run()
}

func main() {
	Run()
}
