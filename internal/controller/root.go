package controller

import (
	"meditationbe/internal/service"

	"go.uber.org/zap"
)

type RootController struct {
	log *zap.Logger
	userService service.UserService
	audioService service.AudioService
} 

func NewRootController(userService service.UserService, audioService service.AudioService, log *zap.Logger) *RootController {
	return &RootController{
		userService: userService,
		audioService: audioService,
		log: log,
	}
}