package controller

import (
	"context"
	"meditationbe/internal/dto"
	"meditationbe/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid/v5"
)

// Login godoc
//
//	@Summary	Login
//	@Tags		user
//	@Produce	json
//	@Param		login	body		dto.UserLoginPayload	true	"Login user"
//	@Success	200		{object}	dto.LoginResponse
//	@Failure	400		{string}	string
//	@Router		/login [post]
//	@deprecated
// func (r *RootController) Login(c *fiber.Ctx) error {
// 	payload := dto.UserLoginPayload{}
// 	if err := c.BodyParser(&payload); err != nil {
// 		return c.SendStatus(fiber.StatusBadRequest)
// 	}

// 	token, err := r.userService.Login(context.Background(), &payload)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
// 	}

// 	return c.JSON(dto.LoginResponse{Token: token})
// }

// Register godoc
//
//	@Summary	Register
//	@Tags		user
//	@Param		register	body		dto.UserRegisterPayload	true	"Register user"
//	@Success	200			{string}	string
//	@Failure	400			{string}	string
//	@Router		/register [post]
//	@deprecated
// func (r *RootController) Register(c *fiber.Ctx) error {
// 	payload := dto.UserRegisterPayload{}
// 	if err := c.BodyParser(&payload); err != nil {
// 		return c.SendStatus(fiber.StatusBadRequest)
// 	}

// 	r.log.Debug(payload.Email)

// 	err := r.userService.Register(context.Background(), &payload)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
// 	}

// 	return c.SendStatus(fiber.StatusOK)
// }

// GetUser godoc
//
//	@Summary	GetUser
//	@Tags		user
//	@Success	200	{object}	dto.UserResponse
//	@Failure	400	{string}	string
//	@Failure	401	{string}	string
//	@Router		/me [get]
//	@Security	BearerToken
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

	user, err := r.userService.GetByID(context.Background(), userUuid)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("no user found")
	}

	userResponse := dto.UserResponse{
		ID:        user.ID,
		TgID:      user.TgID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		PhotoUrl:  user.PhotoUrl,
		Role:      user.Role,
	}

	return c.JSON(userResponse)
}
