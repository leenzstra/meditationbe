package service

import (
	"context"
	"errors"
	"time"

	"meditationbe/config"
	"meditationbe/internal/domain"
	"meditationbe/internal/repository"
	tgauth "meditationbe/internal/tg_auth"
	"meditationbe/internal/utils"

	uuid "github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
)

type UserService interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	Auth(ctx context.Context, creds *tgauth.Credentials) (string, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func (u *userService) Auth(ctx context.Context, creds *tgauth.Credentials) (string, error) {
	if err := creds.Verify([]byte(config.GetConfig().BotToken)); err != nil {
		return "", err
	}

	foundUser, err := u.userRepo.GetByTgID(ctx, creds.ID)

	switch {
	case errors.Is(err, repository.ErrNotFound):
		err = u.register(ctx, creds)
		if err != nil {
			return "", err
		}

		foundUser, err = u.userRepo.GetByTgID(ctx, creds.ID)
		if err != nil {
			return "", err
		}
	case err != nil:
		return "", err
	}

	claims := jwt.MapClaims{
		"sub":   foundUser.ID.String(),
		"admin": foundUser.Role == "admin",
		"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(),
	}

	return utils.GenerateJWT(claims, []byte(config.GetConfig().JWTSecret))
}

func (u *userService) register(ctx context.Context, creds *tgauth.Credentials) error {
	userUuid, err := uuid.NewV4()
	if err != nil {
		return err
	}

	newUser := &domain.User{
		ID:        userUuid,
		TgID:      creds.ID,
		Username:  creds.Username,
		FirstName: creds.FirstName,
		LastName:  creds.LastName,
		PhotoUrl:  creds.PhotoURL,
		Provider:  "telegram",
		Role:      "user",
	}

	return u.userRepo.Add(ctx, newUser)
}

func (u *userService) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	return u.userRepo.GetByUsername(ctx, username)
}

func (u *userService) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return u.userRepo.GetByID(ctx, id)
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		userRepo: repo,
	}
}
