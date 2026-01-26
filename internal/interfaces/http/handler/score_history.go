package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SubmitScoreRequest struct {
	GameID string `json:"game_id"`
	Score  int    `json:"score"`
}

func (h *Handler) submitScore(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	var req SubmitScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	gameID, err := uuid.Parse(req.GameID)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.SubmitScore(c.Request.Context(), userID, gameID, req.Score); err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
