package server

import (
	"meditationbe/internal/controller"
	"meditationbe/internal/database"
	"meditationbe/internal/middleware"
	"meditationbe/internal/repository"
	"meditationbe/internal/service"

	"github.com/gofiber/fiber/v2"
	flog "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
)

func NewRouter(db *database.Database, logger *zap.Logger) *fiber.App {
	router := fiber.New(fiber.Config{
		BodyLimit: 25 * 1024 * 1024, 
	})

	router.Use(flog.New())
	router.Use(recover.New())

	setupRoutes(router, db, logger)

	return router
}

func setupRoutes(router *fiber.App, db *database.Database, logger *zap.Logger) {
	// user related
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	// audio related
	audioRepo := repository.NewAudioRepository(db)
	audioUploader := service.NewServerAudioUploader("./audio")
	audioService := service.NewAudioService(audioRepo, audioUploader)

	// auth related
	authRequired := middleware.AuthRequired()

	// main controller
	root := controller.NewRootController(userService, audioService, logger)

	// routes
	api := router.Group("/api")

	api.Post("/login", root.Login)
	api.Post("/register", root.Register)
	api.Get("/health", root.Status) 
	
	api.Get("/me", authRequired, root.GetUser)

	audioRestricted := api.Group("/audio", authRequired)
	audioRestricted.Post("/upload", root.UploadAudio)
	audioRestricted.Delete("/delete", root.DeleteAudio)
	audioRestricted.Get("/list", root.GetAudioList) 
	audioRestricted.Get("/:uuid", root.GetAudio)
	audioRestricted.Post("/update", root.UpdateAudio)
}
