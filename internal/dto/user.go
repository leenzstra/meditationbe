package dto

import (
	"github.com/gofrs/uuid/v5"
)

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	TgID      int64     `json:"tg_id"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	PhotoUrl  string    `json:"photo_url"`
	Role      string    `json:"role"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
