package controller

import (
	"context"
	"meditationbe/internal/domain"
	"meditationbe/internal/dto"

	"github.com/gofiber/fiber/v2"
)

func (r *RootController) UploadAudio(c *fiber.Ctx) error {
	payload := dto.AudioAddPayload{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err = r.audioService.Add(context.Background(), &payload, file)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

func (r *RootController) DeleteAudio(c *fiber.Ctx) error {
	payload := dto.AudioDeletePayload{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err := r.audioService.Delete(context.Background(), &payload)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

func (r *RootController) GetAudioList(c *fiber.Ctx) error {
	audio, err := r.audioService.GetAll(context.Background())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.JSON(audio)
}

func (r *RootController) GetAudio(c *fiber.Ctx) error {
	payload := dto.AudioQueryPayload{}
	if err := c.ParamsParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	r.log.Debug(payload.UUID.String()) 

	audio, err := r.audioService.Get(context.Background(), payload.UUID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.JSON(audio)
}

func (r *RootController) UpdateAudio(c *fiber.Ctx) error {
	payload := &domain.Audio{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err := r.audioService.Update(context.Background(), payload)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}