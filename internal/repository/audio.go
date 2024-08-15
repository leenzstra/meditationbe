package repository

import (
	"context"
	"fmt"
	"meditationbe/internal/database"
	"meditationbe/internal/domain"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gofrs/uuid/v5"
)

type AudioRepository interface {
	Get(ctx context.Context, uuid uuid.UUID) (*domain.Audio, error)
	GetAll(ctx context.Context) ([]*domain.Audio, error)
	Delete(ctx context.Context, uuid uuid.UUID) error
	Update(ctx context.Context, audio *domain.Audio) error
	Add(ctx context.Context, audio *domain.Audio) error
}

type audioRepository struct {
	db *database.Database
}

// Add implements AudioRepository.
func (r *audioRepository) Add(ctx context.Context, audio *domain.Audio) error {
	sql, args, err := r.db.Builder.
		Insert("audio").Columns("uuid", "name", "description", "path").
		Values(audio.UUID.String(), audio.Name, audio.Description, audio.Path).ToSql()
	if err != nil { 
		return err 
	} 

	r.db.Logger.Debug(fmt.Sprintf("%v %v", sql, args)) 

	_, err = r.db.Pool.Exec(ctx, sql, args...)
	return err
}

// Delete implements AudioRepository.
func (r *audioRepository) Delete(ctx context.Context, uuid uuid.UUID) error {
	sql, args, err := r.db.Builder.
		Delete("audio").Where(sq.Eq{"uuid": uuid}).ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Pool.Exec(ctx, sql, args...)
	return err
}

// Get implements AudioRepository.
func (r *audioRepository) Get(ctx context.Context, uuid uuid.UUID) (*domain.Audio, error) {
	sql, args, err := r.db.Builder.
		Select("*").From("audio").Where(sq.Eq{"uuid": uuid}).ToSql()
	if err != nil {
		return nil, err
	}

	audio := []*domain.Audio{}
	err = pgxscan.Select(ctx, r.db.Pool, &audio, sql, args...)
	if err != nil {
		return nil, err
	}

	if len(audio) == 0 {
		return nil, ErrNotFound
	}

	return audio[0], nil
}

// GetAll implements AudioRepository.
func (r *audioRepository) GetAll(ctx context.Context) ([]*domain.Audio, error) {
	sql, args, err := r.db.Builder.
		Select("*").From("audio").ToSql()
	if err != nil {
		return nil, err
	}

	audio := []*domain.Audio{}
	err = pgxscan.Select(ctx, r.db.Pool, &audio, sql, args...)
	if err != nil {
		return nil, err
	}

	if len(audio) == 0 {
		return nil, ErrNotFound
	}

	return audio, nil
}

// Update implements AudioRepository.
func (r *audioRepository) Update(ctx context.Context, audio *domain.Audio) error {
	setValues := map[string]interface{}{
		"name": audio.Name,
		"description": audio.Description,
		"path": audio.Path,
	}
	sql, args, err := r.db.Builder.
		Update("audio").SetMap(setValues).Where(sq.Eq{"uuid": audio.UUID.String()}).ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Pool.Exec(ctx, sql, args...)
	return err
}

func NewAudioRepository(db *database.Database) AudioRepository { return &audioRepository{db: db} }
