package handler

import (
	"github.com/gin-gonic/gin"
	"log/slog"
)

// Error represents an API error response
// @Description API error information
type Error struct {
	Message string `json:"message"`
}

// StatusResponse represents a simple status response
// @Description Basic status response
type statusResponse struct {
	Status string `json:"status"`
}

// StatusFloat represents a numeric status response
// @Description Numeric status response
type statusFloat struct {
	Status float64 `json:"status"`
}

// NewErrorResponse logs error and returns error response
func NewErrorResponse(c *gin.Context, statusCode int, message string) {
	slog.Error(message)
	c.AbortWithStatusJSON(statusCode, Error{message})
}
