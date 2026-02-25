package handler

import (
	"OnlineLeadership/internal/infrastructure/logger"
	"OnlineLeadership/internal/interfaces/http/middleware"
	"OnlineLeadership/internal/usecase"
	"github.com/gin-contrib/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

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
	r.Use(middleware.PrometheusMiddleware())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(files.Handler))
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Authorization", "Content-Type"},
	}))
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
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
