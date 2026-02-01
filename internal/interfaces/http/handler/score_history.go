package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SubmitScoreInput represents score submission payload
type SubmitScoreInput struct {
	GameID string `json:"game_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
	Score  int    `json:"score" binding:"required,min=0" example:"12345"`
}

// @Summary Submit player's score
// @Description Submit a user's score for a game
// @Tags score
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param input body SubmitScoreInput true "Score info"
// @Success 200 {object} StatusResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/score/submit [post]
func (h *Handler) submitScore(c *gin.Context) {
	userID, err := getUserId(c)
	ctx := c.Request.Context()
	if err != nil {
		h.log.Error(ctx, "get user id error", err.Error())
		NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	var req SubmitScoreInput
	if err := c.ShouldBindJSON(&req); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	gameID, err := uuid.Parse(req.GameID)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid game_id format")
		return
	}

	if err := h.service.SubmitScore(c.Request.Context(), userID, gameID, req.Score); err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, StatusResponse{
		Status: "ok",
	})
}
