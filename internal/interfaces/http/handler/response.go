package handler

import (
	"github.com/gin-gonic/gin"
	"log/slog"
)

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Message string `json:"message" example:"internal server error"`
}

// StatusResponse represents a simple status response
type StatusResponse struct {
	Status string `json:"status" example:"ok"`
}

// RankResponse represents user rank response
type RankResponse struct {
	Rank int64 `json:"rank" example:"1"`
}

// LeaderboardUserDTO represents a user entry in the leaderboard
type LeaderboardUserDTO struct {
	UserID string `json:"user_id" example:"01234567-89ab-cdef-0123-456789abcdef"`
	Score  int64  `json:"score" example:"12345"`
	Rank   int64  `json:"rank" example:"1"`
}

// LeaderboardResponse represents leaderboard list response
type LeaderboardResponse struct {
	Data []LeaderboardUserDTO `json:"data"`
}

// GameDTO represents game information
type GameDTO struct {
	ID   string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name string `json:"name" example:"Chess"`
}

// GamesResponse represents games list response
type GamesResponse struct {
	Data []GameDTO `json:"data"`
}

// RegisterResponse represents registration response
type RegisterResponse struct {
	UserID string `json:"user_id" example:"01234567-89ab-cdef-0123-456789abcdef"`
}

// LoginResponse represents login response
type LoginResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// GameIDResponse represents game creation response
type GameIDResponse struct {
	GameID string `json:"game_id" example:"123e4567-e89b-12d3-a456-426614174000"`
}

func NewErrorResponse(c *gin.Context, statusCode int, message string) {
	slog.Error(message)
	c.AbortWithStatusJSON(statusCode, ErrorResponse{Message: message})
}
