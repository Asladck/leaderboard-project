package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "UserId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		NewErrorResponse(c, http.StatusBadRequest, "empty auth head")
		return
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		NewErrorResponse(c, http.StatusBadRequest, "invalid auth header")
		return
	}
	userId, err := h.service.Auth.ParseAccessToken(headerParts[1])
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.Set(userCtx, userId)
}

var ErrUserNotAuthorized = errors.New("user not authorized")

func getUserId(c *gin.Context) (uuid.UUID, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return uuid.UUID{}, ErrUserNotAuthorized
	}

	userID, ok := id.(uuid.UUID)
	if !ok {
		return uuid.UUID{}, ErrUserNotAuthorized
	}

	return userID, nil
}
