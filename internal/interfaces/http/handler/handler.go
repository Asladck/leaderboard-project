package handler

import (
	"OnlineLeadership/internal/infrastructure/logger"
	"OnlineLeadership/internal/usecase"

	"github.com/gin-gonic/gin"
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

	// ===== Auth group =====
	auth := r.Group("/auth")
	{
		auth.POST("/register", h.signUp)
		auth.POST("/login", h.signIn)
	}

	//===== Score group ===== //apply auth middleware as needed, e.g.:
	api := r.Group("/api", h.userIdentity)
	{
		score := api.Group("/score")
		{
			score.POST("/submit", h.submitScore)
		}
		//leaderboard := api.Group("/api/leaderboard")
		//{
		//	leaderboard.GET("/global", func(c *gin.Context) { h.globalLeaderboard(c.Writer, c.Request) })
		//	leaderboard.GET("/me", func(c *gin.Context) { h.myRank(c.Writer, c.Request) })
		//	leaderboard.GET("/top", func(c *gin.Context) { h.topPlayers(c.Writer, c.Request) })
		//}
	}

	return r
}
