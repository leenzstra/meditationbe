package service

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"

	uuid "github.com/gofrs/uuid/v5"
	"github.com/supabase-community/storage-go"
)

type AudioUploader interface {
	Upload(ctx context.Context, uuid uuid.UUID, file *bytes.Buffer) (string, error)
	Delete(ctx context.Context, uuid uuid.UUID) error
}

type filesystemAudioUploader struct {
	basePath string
}

// Delete implements AudioUploader.
func (u *filesystemAudioUploader) Delete(ctx context.Context, uuid uuid.UUID) error {
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
func (u *filesystemAudioUploader) Upload(ctx context.Context, uuid uuid.UUID, file *bytes.Buffer) (string, error) {
	filepath := filepath.Join(u.basePath, fmt.Sprint(uuid.String(), ".mp3"))
	if filepath == "" {
		return "", fmt.Errorf("audio delete error: empty path")
	}

	u.ensureDir(u.basePath)
	if u.fileExists(filepath) {
		return "", fmt.Errorf("audio delete error: file exists")
	}

	if err := os.WriteFile(filepath, file.Bytes(), 0644); err != nil {
		return "", err
	}

	return filepath, nil
}

func (u *filesystemAudioUploader) fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (u *filesystemAudioUploader) ensureDir(dir string) error {
	err := os.MkdirAll(dir, os.ModeDir)

	if err == nil || os.IsExist(err) {
		return nil
	} else {
		return err
	}
}

func NewServerAudioUploader(basePath string) AudioUploader {
	return &filesystemAudioUploader{
		basePath: basePath,
	}
}

type s3AudioUploader struct {
	storage *storage_go.Client
	bucketId string
}

// Delete implements AudioUploader.
func (s *s3AudioUploader) Delete(ctx context.Context, uuid uuid.UUID) error {
	filePath := fmt.Sprintf("%s.mp3", uuid.String())
	removeResp, err := s.storage.RemoveFile(s.bucketId, []string{filePath})
	if err != nil {
		return err 
	}

	if removeResp[0].Error != "" {
		return fmt.Errorf(removeResp[0].Error)
	}

	return nil
}

// Upload implements AudioUploader.
func (s *s3AudioUploader) Upload(ctx context.Context, uuid uuid.UUID, file *bytes.Buffer) (string, error) {
	filePath := fmt.Sprintf("%s.mp3", uuid.String()) 
	contentType := "audio/mpeg"

	uploadResp, err := s.storage.UploadFile(s.bucketId, filePath, file, storage_go.FileOptions{
		ContentType: &contentType,
	})
	if err != nil {
		return "", fmt.Errorf("upload: %s", err) 
	}

	if uploadResp.Error != "" {
		return "", fmt.Errorf("upload: %s", uploadResp.Error)
	}

	publicResp := s.storage.GetPublicUrl(s.bucketId, filePath)

	if publicResp.SignedURL == "" {
		return "", fmt.Errorf("upload: no public url")
	}

	return publicResp.SignedURL , nil
}

func NewS3AudioUploader(storage *storage_go.Client, bucketId string) AudioUploader {
	return &s3AudioUploader{
		storage: storage,
		bucketId: bucketId,
	}
}
