package dto

import "github.com/gofrs/uuid/v5"

type UserRegisterPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	UUID  uuid.UUID `json:"uuid"`
	Email string    `json:"email"`
	Role  string    `json:"role"`
}
