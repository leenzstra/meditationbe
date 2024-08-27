package controller

import (
	"context"
	"errors"
	"meditationbe/internal/dto"
	tgauth "meditationbe/internal/tg_auth"

	"github.com/gofiber/fiber/v2"
)

// Telegram Authentication godoc
//
//	@Summary	Telegram auth
//	@Tags		user
//	@Produce	json
//	@Param		login	query		tgauth.Credentials	true	"user creds"
//	@Success	200		{object}	dto.LoginResponse
//	@Failure	400		{string}	string
//	@Router		/auth/telegram [get]
func (r *RootController) TelegramAuth(c *fiber.Ctx) error {
	payload := tgauth.Credentials{}
	if err := c.QueryParser(&payload); err != nil && !errors.Is(err, fiber.ErrUnprocessableEntity) {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	token, err := r.userService.Auth(context.Background(), &payload)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.JSON(dto.LoginResponse{Token: token})
}