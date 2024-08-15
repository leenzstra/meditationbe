package controller

import (
	"context"
	"meditationbe/internal/dto"
	"meditationbe/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid/v5"
)

func (r *RootController) Login(c *fiber.Ctx) error {
	payload := dto.UserLoginPayload{}
	if err := c.BodyParser(&payload); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
 
	r.log.Debug(payload.Email)

	token, err := r.userService.Login(context.Background(), &payload)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.JSON(fiber.Map{"token": token})
}

func (r *RootController) Register(c *fiber.Ctx) error {
	payload := dto.UserRegisterPayload{}
	if err := c.BodyParser(&payload); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
 
	r.log.Debug(payload.Email)

	err := r.userService.Register(context.Background(), &payload)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

func (r *RootController) GetUser(c *fiber.Ctx) error {
	token, err := utils.TokenFromLocals(c)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	sub, err := token.Claims.GetSubject()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("not subject found")
	}

	userUuid, err := uuid.FromString(sub)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("uuid parse error")
	}

	user, err := r.userService.GetByUUID(context.Background(), userUuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("no user found")
	}

	userResponse := dto.UserResponse{
		UUID:  user.UUID,
		Email: user.Email,
		Role:  user.Role,
	}

	return c.JSON(userResponse)
}
