package handler

import (
	"OnlineLeadership/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

// RegisterInput represents user registration payload
type RegisterInput struct {
	Username string `json:"username" binding:"required" example:"john_doe"`
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
}

// LoginInput represents user login payload
type LoginInput struct {
	Username string `json:"username" binding:"required" example:"john_doe"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// @Summary Register new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param input body RegisterInput true "Register input"
// @Success 201 {object} RegisterResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/register [post]
func (h *Handler) signUp(c *gin.Context) {
	ctx := c.Request.Context()

	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Service returns uuid.UUID
	userID, err := h.service.Register(ctx, domain.User{
		Username: input.Username,
		Password: input.Password,
		Email:    input.Email,
	})
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Convert uuid.UUID to string for DTO
	c.JSON(http.StatusCreated, RegisterResponse{
		UserID: userID.String(),
	})
}

// @Summary Login user
// @Description Authenticate user and return access and refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param input body LoginInput true "Login input"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/login [post]
func (h *Handler) signIn(c *gin.Context) {
	ctx := c.Request.Context()

	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	at, rt, err := h.service.Login(ctx, input.Username, input.Password)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		AccessToken:  at,
		RefreshToken: rt,
	})
}
