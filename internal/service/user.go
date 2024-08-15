package service

import (
	"context"
	"fmt"
	"time"

	"meditationbe/config"
	"meditationbe/internal/domain"
	"meditationbe/internal/dto"
	"meditationbe/internal/repository"

	uuid "github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Get(ctx context.Context, email string) (*domain.User, error)
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.User, error)
	Add(ctx context.Context, user *domain.User) error
	Register(ctx context.Context, user *dto.UserRegisterPayload) error
	Login(ctx context.Context, user *dto.UserLoginPayload) (string, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func (u *userService) Add(ctx context.Context, user *domain.User) error {
	return u.userRepo.Add(ctx, user)
}

func (u *userService) Login(ctx context.Context, user *dto.UserLoginPayload) (string, error) {
	foundUser, err := u.Get(ctx, user.Email)
	if err != nil {
		return "", err 
	}

	if !u.checkPasswordHash(user.Password, foundUser.PassHash) {
		return "", fmt.Errorf("invalud password")
	}

	// Create claims
	claims := jwt.MapClaims{
		"sub":   foundUser.UUID.String(),
		"email": foundUser.Email,
		"admin": foundUser.Role == "admin",
		"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token
	tokenString, err := token.SignedString([]byte(config.GetConfig().JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (u *userService) Register(ctx context.Context, user *dto.UserRegisterPayload) error {
	userUuid, err := uuid.NewV4()
	if err != nil {
		return err
	}

	passHash, err := u.hashPassword(user.Password)
	if err != nil {
		return err
	}

	newUser := &domain.User{
		UUID:     userUuid,
		Email:    user.Email,
		PassHash: string(passHash),
		Role:     "user",
	}

	return u.Add(ctx, newUser)

}

func (u *userService) Get(ctx context.Context, email string) (*domain.User, error) {
	return u.userRepo.Get(ctx, email)
}

func (u *userService) GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.User, error) {
	return u.userRepo.GetByUUID(ctx, uuid)
}

func (u *userService) hashPassword(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), 14)
	return string(bytes), err
}

func (u *userService) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		userRepo: repo,
	}
}
