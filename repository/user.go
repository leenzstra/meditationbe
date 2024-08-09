package repository

import (
	"context"

	"github.com/leenzstra/meditationbe/db"
	"github.com/leenzstra/meditationbe/domain"
)

type IUserRepository interface {
	SignUp(ctx context.Context, user *domain.User) (string, error)
	SignIn(ctx context.Context, user *domain.User) (string, error)
	Get(ctx context.Context, uid string) (*domain.User, error)
	Delete(ctx context.Context, uid string) (bool, error)
}

type userRepo struct {
	db *db.Database
}

// Delete implements IUserRepository.
func (u *userRepo) Delete(ctx context.Context, uid string) (bool, error) {
	panic("unimplemented")
}

// Get implements IUserRepository.
func (u *userRepo) Get(ctx context.Context, uid string) (*domain.User, error) {
	panic("unimplemented")
}

// Login implements IUserRepository.
func (u *userRepo) Login(ctx context.Context, user *domain.User) (string, error) {
	panic("unimplemented")
}

// Register implements IUserRepository.
func (u *userRepo) Register(ctx context.Context, user *domain.User) (string, error) {
	panic("unimplemented")
}

func NewUserRepository(db *db.Database) IUserRepository {
	return &userRepo{db: db}
}
