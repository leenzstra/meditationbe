package domain

import "github.com/gofrs/uuid/v5"

type User struct {
	ID        uuid.UUID `json:"id"`
	TgID      int64     `json:"tg_id"`
	Username  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	PhotoUrl  string    `json:"photo_url"`
	Provider  string    `json:"provider"`
	Role      string    `json:"role"`
}
