package domain

import "github.com/google/uuid"

// Game представляет игру.
// swagger:model Game
type Game struct {
	Id   uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name string    `json:"name" example:"Chess"`
}
