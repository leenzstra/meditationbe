package domain

import "github.com/gofrs/uuid/v5"

type User struct {
	UUID     uuid.UUID `json:"uuid"`
	Email    string `json:"email"`
	PassHash string `json:"pass_hash"`
	Role     string `json:"role"`
}