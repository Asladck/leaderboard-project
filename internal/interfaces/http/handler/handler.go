package handler

import (
	"OnlineLeadership/internal/infrastructure/logger"
	"OnlineLeadership/internal/usecase"

	"github.com/gin-gonic/gin"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	service *usecase.Service
	log     *logger.SlogLogger
}

func NewHandler(service *usecase.Service, log *logger.SlogLogger) *Handler {
	return &Handler{service: service, log: log}
}

func (h *Handler) InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))

	// Auth endpoints
	auth := r.Group("/auth")
	{
		auth.POST("/register", h.signUp)
		auth.POST("/login", h.signIn)
	}

	// Admin endpoints
	admin := r.Group("/admin")
	{
		admin.POST("/create", h.createGame)
		admin.GET("/games", h.getGames)
	}

	// Protected API endpoints
	api := r.Group("/api", h.userIdentity)
	{
		score := api.Group("/score")
		{
			score.POST("/submit", h.submitScore)
		}
		leaderboard := api.Group("/leaderboard")
		{
			leaderboard.GET("/global", h.globalLeaderboard)
			leaderboard.GET("/my", h.myRank)
			leaderboard.POST("/top", h.topPlayers)
		}
	}

	return r
}
