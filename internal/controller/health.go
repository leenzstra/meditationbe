package controller

import (
	"github.com/gofiber/fiber/v2"
)

// Status godoc
//	@Summary		API ststus
//	@Description	get API status
//	@Tags			health
//	@Produce		json
//	@Success		200	{string}	string
//	@Router			/health [get]
func (r *RootController) Status(c *fiber.Ctx) error {
	return c.SendString("ok")
}
