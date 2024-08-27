package service

import (
	"bytes"
	"context"
	"meditationbe/internal/domain"
	"meditationbe/internal/dto"
	"meditationbe/internal/repository"

	uuid "github.com/gofrs/uuid/v5"
)

type AudioService interface {
	Get(ctx context.Context, id uuid.UUID) (*domain.Audio, error)
	GetAll(ctx context.Context) ([]*domain.Audio, error)
	Delete(ctx context.Context, audio *dto.AudioDeletePayload) error
	Update(ctx context.Context, audio *domain.Audio) error
	Add(ctx context.Context, audio *dto.AudioAddPayload, file *bytes.Buffer) error
}

type audioService struct {
	audioRepo repository.AudioRepository
	uploader  AudioUploader
}

// Add implements AudioService.
func (s *audioService) Add(ctx context.Context, audio *dto.AudioAddPayload, file *bytes.Buffer) error {
	audioId, err := uuid.NewV4()
	if err != nil {
		return err
	}

	path, err := s.uploader.Upload(ctx, audioId, file)
	if err != nil {
		return err
	}

	newAudio := &domain.Audio{
		ID:          audioId,
		Name:        audio.Name,
		Description: audio.Description,
		Path:        path,
	}

	err = s.audioRepo.Add(ctx, newAudio)
	if err != nil {
		return err
	}

	return nil
}

// Delete implements AudioService.
func (s *audioService) Delete(ctx context.Context, audio *dto.AudioDeletePayload) error {
	err := s.uploader.Delete(ctx, audio.ID)
	if err != nil {
		return err
	}

	err = s.audioRepo.Delete(ctx, audio.ID)
	if err != nil {
		return err
	}

	return nil
}

// Get implements AudioService.
func (s *audioService) Get(ctx context.Context, id uuid.UUID) (*domain.Audio, error) {
	return s.audioRepo.Get(ctx, id)
}

// GetAll implements AudioService.
func (s *audioService) GetAll(ctx context.Context) ([]*domain.Audio, error) {
	return s.audioRepo.GetAll(ctx)
}

// Update implements AudioService.
func (s *audioService) Update(ctx context.Context, audio *domain.Audio) error {
	return s.audioRepo.Update(ctx, audio)
}

func NewAudioService(audioRepo repository.AudioRepository, uploader AudioUploader) AudioService {
	return &audioService{
		audioRepo: audioRepo,
		uploader:  uploader,
	}
}
