package controller

import "github.com/gin-gonic/gin"

type HealthController struct{}

func NewHealthController() *HealthController { return &HealthController{} }

func (h *HealthController) Status(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
