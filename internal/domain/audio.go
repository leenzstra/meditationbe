package domain

import "github.com/gofrs/uuid/v5"

type Audio struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Path        string    `json:"path"`
	Owner       uuid.UUID `json:"owner"`
}
