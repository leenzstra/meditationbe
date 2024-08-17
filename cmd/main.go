package main

import (
	"meditationbe/config"
	"meditationbe/internal/database"
	"meditationbe/internal/server"
	"go.uber.org/zap"
	"log"
	_ "meditationbe/docs"
)

func Run() {
	var err error

	config.Init("development")
    cfg := config.GetConfig()

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}

	db, err := database.NewPostgres(cfg.PgUrl, logger)
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	server.NewRouter(db, logger).Listen(":8080")
}

//	@title						Meditation API
//	@version					0.0.1
//	@description				Meditation API spec
//	@host						localhost:8080
//	@BasePath					/api
//	@securityDefinitions.apikey	BearerToken
//	@in							header 
//	@name						Authorization
func main() {
	Run()
}
