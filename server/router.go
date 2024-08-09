package server

import (
	"github.com/gin-gonic/gin"
	"github.com/leenzstra/meditationbe/controller"
	"github.com/leenzstra/meditationbe/db"
	"go.uber.org/zap"
)

func NewRouter(db *db.Database, logger *zap.Logger) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	setupRoutes(router, db, logger)

	return router
}

func setupRoutes(router *gin.Engine, db *db.Database, logger *zap.Logger) {
	health := controller.NewHealthController()

	router.GET("/health", health.Status)
}
