package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateGameInput represents input for creating a game
type CreateGameInput struct {
	Name string `json:"name" binding:"required" example:"Chess"`
}

// @Summary Create a new game
// @Description Create a new game with given name
// @Tags admin
// @Accept json
// @Produce json
// @Param input body CreateGameInput true "Game input"
// @Success 201 {object} GameIDResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/create [post]
func (h *Handler) createGame(c *gin.Context) {
	ctx := c.Request.Context()
	var input CreateGameInput
	if err := c.ShouldBindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Service returns uuid.UUID
	id, err := h.service.Admin.Create(ctx, input.Name)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Convert uuid.UUID to string for DTO
	c.JSON(http.StatusCreated, GameIDResponse{
		GameID: id.String(),
	})
}

// @Summary Get list of games
// @Description Returns all games
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} GamesResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/games [get]
func (h *Handler) getGames(c *gin.Context) {
	ctx := c.Request.Context()
	games, err := h.service.Admin.GetGames(ctx)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Convert domain models to DTOs (uuid.UUID -> string)
	gameDTOs := make([]GameDTO, 0, len(games))
	for _, game := range games {
		gameDTOs = append(gameDTOs, GameDTO{
			ID:   game.Id.String(),
			Name: game.Name,
		})
	}

	c.JSON(http.StatusOK, GamesResponse{
		Data: gameDTOs,
	})
}
