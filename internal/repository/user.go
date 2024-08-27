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

var (
	ErrNotFound = fmt.Errorf("not found")
)

type UserRepository interface {
	Add(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByTgID(ctx context.Context, tgId int64) (*domain.User, error)
	Delete(ctx context.Context, id string) error
}

type userRepo struct {
	db *database.Database
}

func (u *userRepo) Add(ctx context.Context, user *domain.User) error {
	sql, args, err := u.db.Builder.
		Insert("users").Columns("id", "tg_id", "username", "first_name", "last_name", "photo_url", "provider", "role").
		Values(user.ID.String(), user.TgID, user.Username, user.FirstName, user.LastName, user.PhotoUrl, user.Provider, user.Role).ToSql()
	if err != nil {
		return err
	}

	u.db.Logger.Debug(fmt.Sprintf("%v %v", sql, args))

	_, err = u.db.Pool.Exec(ctx, sql, args...)
	return err
}

func (u *userRepo) Update(ctx context.Context, user *domain.User) error {
	setValues := map[string]interface{}{
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"photo_url":  user.PhotoUrl,
	}
	sql, args, err := u.db.Builder.
		Update("users").SetMap(setValues).Where(sq.Eq{"id": user.ID.String()}).ToSql()
	if err != nil {
		return err
	}

	_, err = u.db.Pool.Exec(ctx, sql, args...)
	return err
}

func (u *userRepo) Delete(ctx context.Context, id string) error {
	sql, args, err := u.db.Builder.
		Delete("users").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return err
	}

	_, err = u.db.Pool.Exec(ctx, sql, args...)
	return err
}

func (u *userRepo) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	return u.getBy(ctx, "username", username)
}

func (u *userRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return u.getBy(ctx, "id", id.String())
}

func (u *userRepo) GetByTgID(ctx context.Context, tgid int64) (*domain.User, error) {
	return u.getBy(ctx, "tg_id", tgid)
}

func (u *userRepo) getBy(ctx context.Context, field string, value any) (*domain.User, error) {
	sql, args, err := u.db.Builder.
		Select("*").From("users").Where(sq.Eq{field: value}).ToSql()
	if err != nil {
		return nil, err
	}

	users := []*domain.User{}
	err = pgxscan.Select(ctx, u.db.Pool, &users, sql, args...)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, ErrNotFound
	}

	return users[0], nil
}

func NewUserRepository(db *database.Database) UserRepository {
	return &userRepo{db: db}
}
