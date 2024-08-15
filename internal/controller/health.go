package controller

import (
	"meditationbe/internal/utils"

	"github.com/gofiber/fiber/v2"
)

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

	return c.JSON(fiber.Map{
		"status": "ok",
		"user":   user,
	})
}
