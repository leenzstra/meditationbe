package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"meditationbe/internal/domain"
	"meditationbe/internal/dto"
	"meditationbe/internal/repository"
	"os"
	"path/filepath"

	uuid "github.com/gofrs/uuid/v5"
)

type AudioService interface {
	Get(ctx context.Context, uuid uuid.UUID) (*domain.Audio, error)
	GetAll(ctx context.Context) ([]*domain.Audio, error)
	Delete(ctx context.Context, audio *dto.AudioDeletePayload) error
	Update(ctx context.Context, audio *domain.Audio) error
	Add(ctx context.Context, audio *dto.AudioAddPayload, file io.Reader) error
}

type audioService struct {
	audioRepo repository.AudioRepository
	uploader  AudioUploader
}

// Add implements AudioService.
func (s *audioService) Add(ctx context.Context, audio *dto.AudioAddPayload, file io.Reader) error {
	audioUuid, err := uuid.NewV4()
	if err != nil {
		return err
	}

	path, err := s.uploader.Upload(ctx, audioUuid, file)
	if err != nil {
		return err
	}

	newAudio := &domain.Audio{
		UUID:        audioUuid,
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
	err := s.uploader.Delete(ctx, audio.UUID)
	if err != nil {
		return err
	}

	err = s.audioRepo.Delete(ctx, audio.UUID)
	if err != nil {
		return err
	}

	return nil
}

// Get implements AudioService.
func (s *audioService) Get(ctx context.Context, uuid uuid.UUID) (*domain.Audio, error) {
	return s.audioRepo.Get(ctx, uuid)
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

type AudioUploader interface {
	Upload(ctx context.Context, uuid uuid.UUID, file io.Reader) (string, error)
	Delete(ctx context.Context, uuid uuid.UUID) error
}

type serverAudioUploader struct {
	basePath string
}

// Delete implements AudioUploader.
func (u *serverAudioUploader) Delete(ctx context.Context, uuid uuid.UUID) error {
	filepath := filepath.Join(u.basePath, fmt.Sprint(uuid.String(), ".mp3"))
	if filepath == "" {
		return fmt.Errorf("audio delete error: empty path")
	}

	if !u.fileExists(filepath) {
		return fmt.Errorf("audio delete error: file not found")
	}

	if err := os.Remove(filepath); err != nil {
		return err
	}

	return nil
}

// Upload implements AudioUploader.
func (u *serverAudioUploader) Upload(ctx context.Context, uuid uuid.UUID, file io.Reader) (string, error) {
	filepath := filepath.Join(u.basePath, fmt.Sprint(uuid.String(), ".mp3"))
	if filepath == "" {
		return "", fmt.Errorf("audio delete error: empty path")
	}

	u.ensureDir(u.basePath)
	if u.fileExists(filepath) {
		return "", fmt.Errorf("audio delete error: file exists")
	}

	buf := bytes.Buffer{}
	_, err := buf.ReadFrom(file)
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(filepath, buf.Bytes(), 0644); err != nil {
		return "", err
	}

	return filepath, nil
}

func (u *serverAudioUploader) fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (u *serverAudioUploader) ensureDir(dir string) error {
    err := os.MkdirAll(dir, os.ModeDir)

    if err == nil || os.IsExist(err) {
        return nil
    } else {
        return err
    }
}

func NewServerAudioUploader(basePath string) AudioUploader {
	return &serverAudioUploader{
		basePath: basePath,
	}
}
