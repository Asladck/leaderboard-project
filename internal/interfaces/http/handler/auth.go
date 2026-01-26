package handler

import (
	"OnlineLeadership/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RegisterInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) signUp(c *gin.Context) {
	ctx := c.Request.Context()

	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	userID, err := h.service.Register(ctx, domain.User{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 4️⃣ Ответ
	c.JSON(http.StatusCreated, gin.H{
		"user_id": userID,
	})
}

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

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

	// 4️⃣ Ответ
	c.JSON(http.StatusOK, gin.H{
		"access_token":  at,
		"refresh_token": rt,
	})
}
