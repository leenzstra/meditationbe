package dto

import "github.com/gofrs/uuid/v5"

type AudioAddPayload struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AudioDeletePayload struct {
	UUID uuid.UUID `json:"uuid"`
}

type AudioQueryPayload struct {
	UUID uuid.UUID `json:"uuid" params:"uuid"`
}

type AudioUpdatePayload struct {
	UUID        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Path        string    `json:"path"`
}
