package dto

import "github.com/gofrs/uuid/v5"

type AudioAddPayload struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AudioDeletePayload struct {
	ID uuid.UUID `json:"id"`
}

type AudioQueryPayload struct {
	ID uuid.UUID `json:"id" params:"id"`
}

type AudioUpdatePayload struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Path        string    `json:"path"`
}
