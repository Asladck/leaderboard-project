package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"net/http"
	"strconv"
)

// TopPlayersInput represents input for getting top players
type TopPlayersInput struct {
	GameID string `json:"game_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// @Summary Get global leaderboard
// @Description Returns a paginated global leaderboard
// @Tags leaderboard
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param offset query int false "Offset" default(0)
// @Param limit query int false "Limit" default(50) maximum(100)
// @Success 200 {object} LeaderboardResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/leaderboard/global [get]
func (h *Handler) globalLeaderboard(c *gin.Context) {
	limit := 50
	offset := 0
	ctx := c.Request.Context()
	if o := c.Query("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}

	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 100 {
			limit = v
		}
	}
	users, err := h.service.GetGlobalLeaderboard(
		ctx,
		offset,
		limit,
	)
	if err != nil {
		h.log.Error(ctx, "get global leaderboard failed", err)
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Convert domain models to DTOs (uuid.UUID -> string)
	userDTOs := make([]LeaderboardUserDTO, 0, len(users))
	for _, user := range users {
		userDTOs = append(userDTOs, LeaderboardUserDTO{
			UserID: user.UserID.String(),
			Score:  user.Score,
			Rank:   user.Rank,
		})
	}

	c.JSON(http.StatusOK, LeaderboardResponse{
		Data: userDTOs,
	})
}

// @Summary Get current user's rank
// @Description Returns the rank of the authenticated user in the global leaderboard
// @Tags leaderboard
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} RankResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/leaderboard/my [get]
func (h *Handler) myRank(c *gin.Context) {
	ctx := c.Request.Context()
	userID, err := getUserId(c)
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	rank, err := h.service.Leaderboard.GetMyRank(ctx, userID)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, RankResponse{
		Rank: rank,
	})
}

// @Summary Get top players for a game
// @Description Returns the top players for a specified game
// @Tags leaderboard
// @Accept json
// @Security ApiKeyAuth
// @Produce json
// @Param input body TopPlayersInput true "Game id"
// @Success 200 {object} LeaderboardResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/leaderboard/top [post]
func (h *Handler) topPlayers(c *gin.Context) {
	ctx := c.Request.Context()
	var req TopPlayersInput
	err := c.ShouldBindJSON(&req)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Parse string GameID to uuid.UUID
	gameID, err := uuid.Parse(req.GameID)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid game_id format")
		return
	}

	users, err := h.service.Leaderboard.GetLeaderboard(ctx, gameID)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Convert domain models to DTOs (uuid.UUID -> string)
	userDTOs := make([]LeaderboardUserDTO, 0, len(users))
	for _, user := range users {
		userDTOs = append(userDTOs, LeaderboardUserDTO{
			UserID: user.UserID.String(),
			Score:  user.Score,
			Rank:   user.Rank,
		})
	}

	c.JSON(http.StatusOK, LeaderboardResponse{
		Data: userDTOs,
	})
}
