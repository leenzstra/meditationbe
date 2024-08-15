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
	Add(ctx context.Context, user *domain.User) (error)
	Update(ctx context.Context, user *domain.User) (error)
	Get(ctx context.Context, email string) (*domain.User, error)
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.User, error)
	Delete(ctx context.Context, uid string) (error)
}

type userRepo struct {
	db *database.Database
}

func (u *userRepo) Add(ctx context.Context, user *domain.User) (error) {
	sql, args, err := u.db.Builder.
		Insert("users").Columns("uuid", "email", "role", "pass_hash").
		Values(user.UUID.String(), user.Email, user.Role, user.PassHash).ToSql()
	if err != nil { 
		return err 
	} 

	u.db.Logger.Debug(fmt.Sprintf("%v %v", sql, args)) 

	_, err = u.db.Pool.Exec(ctx, sql, args...)
	return err
}

func (u *userRepo) Update(ctx context.Context, user *domain.User) (error) {
	setValues := map[string]interface{}{
		"email": user.Email,
		"role": user.Role,
		"pass_hash": user.PassHash,
	}
	sql, args, err := u.db.Builder.
		Update("users").SetMap(setValues).Where(sq.Eq{"uuid": user.UUID.String()}).ToSql()
	if err != nil {
		return err
	}

	_, err = u.db.Pool.Exec(ctx, sql, args...)
	return err
}

func (u *userRepo) Delete(ctx context.Context, uuid string) (error) {
	sql, args, err := u.db.Builder.
		Delete("users").Where(sq.Eq{"uuid": uuid}).ToSql()
	if err != nil {
		return err
	}

	_, err = u.db.Pool.Exec(ctx, sql, args...)
	return err
}

func (u *userRepo) Get(ctx context.Context, email string) (*domain.User, error) {
	sql, args, err := u.db.Builder.
		Select("*").From("users").Where(sq.Eq{"email": email}).ToSql()
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

func (u *userRepo) GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.User, error) {
	sql, args, err := u.db.Builder.
		Select("*").From("users").Where(sq.Eq{"uuid": uuid.String()}).ToSql()
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
