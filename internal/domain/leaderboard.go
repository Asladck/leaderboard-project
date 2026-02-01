package domain

import "github.com/google/uuid"

// LeaderboardUser swagger:model LeaderboardUser
// LeaderboardUser представляет запись пользователя в таблице лидеров.
type LeaderboardUser struct {
	// UserID is the unique identifier of the user.
	// example: 01234567-89ab-cdef-0123-456789abcdef
	UserID uuid.UUID `json:"user_id" example:"01234567-89ab-cdef-0123-456789abcdef"`
	// Score is the user's score.
	// example: 12345
	Score int64 `json:"score" example:"12345"`
	// Rank is the user's rank on the leaderboard.
	// example: 1
	Rank int64 `json:"rank" example:"1"`
}
