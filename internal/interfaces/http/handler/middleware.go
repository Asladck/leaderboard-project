package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "UserId"
)

// userIdentity is a Gin middleware that extracts the user id from a Bearer access token.
//
// Swagger annotations for documentation generators (e.g., swaggo):
// @Summary Authenticate user by access token (middleware)
// @Description Parses the "Authorization: Bearer {token}" header, validates the access token and stores the user id in the Gin context under key `UserId`.
// @Tags middleware
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Security ApiKeyAuth
func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		NewErrorResponse(c, http.StatusBadRequest, "empty auth header")
		return
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		NewErrorResponse(c, http.StatusBadRequest, "invalid auth header")
		return
	}

	// Service returns uuid.UUID
	userId, err := h.service.Auth.ParseAccessToken(c.Request.Context(), headerParts[1])
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	// Store uuid.UUID in context
	c.Set(userCtx, userId)
	c.Next()
}

var ErrUserNotAuthorized = errors.New("user not authorized")

// getUserId retrieves the user UUID stored in Gin context by the userIdentity middleware.
// Returns ErrUserNotAuthorized when the value is missing or not a UUID.
func getUserId(c *gin.Context) (uuid.UUID, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		log.Println("userCtx not found")
		return uuid.UUID{}, ErrUserNotAuthorized
	}

	log.Printf("userCtx value: %#v, type: %T\n", id, id)

	userID, ok := id.(uuid.UUID)
	if !ok {
		return uuid.UUID{}, ErrUserNotAuthorized
	}

	return userID, nil
}
