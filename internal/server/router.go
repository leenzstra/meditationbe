package server

import (
	"fmt"
	"meditationbe/config"
	"meditationbe/internal/controller"
	"meditationbe/internal/database"
	"meditationbe/internal/middleware"
	"meditationbe/internal/repository"
	"meditationbe/internal/service"

	"github.com/gofiber/fiber/v2"
	flog "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger" // swagger handler
	"github.com/supabase-community/storage-go"
	"go.uber.org/zap"
)

func NewRouter(db *database.Database, logger *zap.Logger) *fiber.App {
	router := fiber.New(fiber.Config{
		BodyLimit: 25 * 1024 * 1024, 
	})

	router.Use(flog.New())
	router.Use(recover.New())
	router.Get("/swagger/*", swagger.HandlerDefault)

	router.Static("/static", "static")

	setupRoutes(router, db, logger)

	return router
}

func setupRoutes(router *fiber.App, db *database.Database, logger *zap.Logger) {
	cfg := config.GetConfig()
	// user related
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	// audio related
	audioRepo := repository.NewAudioRepository(db)

	storage := storage_go.NewClient(fmt.Sprintf("https://%s/storage/v1", cfg.S3Endpoint), cfg.S3JWT, nil)

	audioUploader := service.NewS3AudioUploader(storage, cfg.S3BucketId)
	audioService := service.NewAudioService(audioRepo, audioUploader)

	// auth related
	authRequired := middleware.AuthRequired()

	// main controller
	root := controller.NewRootController(userService, audioService, logger)

	// routes
	api := router.Group("/api")

	api.Post("/auth/telegram", root.TelegramAuth)
	api.Get("/health", root.Status) 
	
	api.Get("/me", authRequired, root.GetUser)

	audioRestricted := api.Group("/audio", authRequired)
	audioRestricted.Post("/upload", root.UploadAudio)
	audioRestricted.Delete("/delete", root.DeleteAudio)
	audioRestricted.Get("/list", root.GetAudioList) 
	audioRestricted.Get("/:uuid", root.GetAudio)
	audioRestricted.Post("/update", root.UpdateAudio)
}
