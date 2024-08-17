package controller

import (
	"context"
	"meditationbe/internal/domain"
	"meditationbe/internal/dto"

	"github.com/gofiber/fiber/v2"
)

// UploadAudio godoc
//
//	@Summary	UploadAudio
//	@Tags		audio
//	@Accept		x-www-form-urlencoded
//	@Param		file	formData	file				true	"Audio file"
//	@Param		data	formData	dto.AudioAddPayload	true	"Audio data"
//	@Success	200		{string}	string
//	@Failure	400		{string}	string
//	@Router		/audio/upload [post]
//	@Security	BearerToken
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
	defer file.Close()

	err = r.audioService.Add(context.Background(), &payload, file)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

// DeleteAudio godoc
//
//	@Summary	DeleteAudio
//	@Tags		audio
//	@Param		data	body		dto.AudioDeletePayload	true	"Payload"
//	@Success	200		{string}	string
//	@Failure	400		{string}	string
//	@Router		/audio/delete [delete]
//	@Security	BearerToken
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

// GetAudioList godoc
//
//	@Summary	GetAudioList
//	@Tags		audio
//	@Success	200	{array}		domain.Audio	
//	@Failure	400	{string}	string
//	@Router		/audio/list [get]
//	@Security	BearerToken
func (r *RootController) GetAudioList(c *fiber.Ctx) error {
	audio, err := r.audioService.GetAll(context.Background())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.JSON(audio)
}

// GetAudio godoc
//
//	@Summary	GetAudio
//	@Tags		audio
//	@Param		uuid	path		string			true	"UUID"
//	@Success	200		{object}	domain.Audio	
//	@Failure	400		{string}	string
//	@Router		/audio/{uuid} [get]
//	@Security	BearerToken
func (r *RootController) GetAudio(c *fiber.Ctx) error {
	payload := dto.AudioQueryPayload{}
	if err := c.ParamsParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	audio, err := r.audioService.Get(context.Background(), payload.UUID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.JSON(audio)
}

// UpdateAudio godoc
//
//	@Summary	UpdateAudio
//	@Tags		audio
//	@Param		data	body		domain.Audio	true	"Data"
//	@Success	200		{string}	string
//	@Failure	400		{string}	string
//	@Router		/audio/update [post]
//	@Security	BearerToken
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