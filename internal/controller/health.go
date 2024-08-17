package controller

import (
	"meditationbe/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type healthResponse struct {
	Status string `json:"status"`
	User   string `json:"user"`
}

// Status godoc
//	@Summary		API ststus
//	@Description	get API status
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	controller.healthResponse
//	@Router			/health [get]
func (r *RootController) Status(c *fiber.Ctx) error {
	var user string

	token, err := utils.TokenFromHeaders(c)
	if err != nil {
		user = "nouser"
	} else {
		subject, err := token.Claims.GetSubject()
		if err != nil {
			user = "nosubject"
		} else {
			user = subject
		}
	}

	return c.JSON(healthResponse{
		Status: "ok",
		User:   user,
	})
}
