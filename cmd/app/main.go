package main

import (
	"OnlineLeadership/internal/infrastructure/auth"
	"OnlineLeadership/internal/infrastructure/logger"
	"OnlineLeadership/internal/infrastructure/postgres"
	"OnlineLeadership/internal/infrastructure/redis"
	"OnlineLeadership/internal/infrastructure/repository"
	"OnlineLeadership/internal/interfaces/http/handler"
	"OnlineLeadership/internal/interfaces/http/middleware"
	"OnlineLeadership/internal/usecase"
	"context"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log := logger.New("dev")
	ctx := context.Background()
	log.Info(ctx, "App is running")
	if err := initConfig(); err != nil {
		log.Error(ctx, "init config error : ", err.Error())
	}
	accessSecret := os.Getenv("JWT_ACCESS_SECRET")
	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")

	if accessSecret == "" || refreshSecret == "" {
		log.Error(ctx, "JWT secrets are not set")
	}

	retryCfg := postgres.RetryConfig{
		MaxAttempts: 10,
		Delay:       3 * time.Second,
		Timeout:     30 * time.Second,
	}
	log.Info(ctx, "db config",
		"sslmode", viper.GetString("db.sslmode"),
	)

	db, err := postgres.ConnectWithRetry(
		ctx,
		retryCfg,
		viper.GetString("db.username"),
		os.Getenv("DB_PASSWORD"),
		viper.GetString("db.host"),
		viper.GetString("db.port"),
		viper.GetString("db.dbname"),
		viper.GetString("db.sslmode"),
	)
	if err != nil {
		log.Error(ctx, "db connect failed", "error", err)
		return
	}
	tokenManager := auth.NewTokenManager(accessSecret, refreshSecret)
	dbredis := redis.InitRedis()
	repos := repository.NewRepository(db, dbredis, log)
	services := usecase.NewService(repos, log, tokenManager)
	handlers := handler.NewHandler(services, log)
	router := handlers.InitRouter()
	routerWithMiddleware := middleware.RequestID(router)
	srv := new(handler.Server)
	go func() {
		log.Info(ctx, "Leaderboard app starting", "port", viper.GetString("port"))
		if err := srv.Run(viper.GetString("port"), routerWithMiddleware); err != nil {
			log.Error(ctx, "server run error", "error", err)
		}
	}()

	log.Info(ctx, "Leaderboard app starting")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	log.Info(ctx, "Leaderboard is shutting down")
	if err := srv.Shutdown(); err != nil {
		log.Error(ctx, "Error occured on server shutting down: ", err.Error())
	}
	if err := db.Close(); err != nil {
		log.Error(ctx, "Error occured on db connection close: ", err.Error())
	}
}

func initConfig() error {
	viper.SetConfigName("config") // config.yml
	viper.SetConfigType("yaml")   // 🔥 важно
	viper.AddConfigPath(".")
	return viper.ReadInConfig()
}
