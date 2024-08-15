package domain

import "github.com/gofrs/uuid/v5"

type Audio struct {
	UUID        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Path        string    `json:"path"`
}
