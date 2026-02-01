package usecase

import (
	"OnlineLeadership/internal/domain"
	"OnlineLeadership/internal/infrastructure/logger"
	"OnlineLeadership/internal/infrastructure/repository"
	"OnlineLeadership/internal/usecase/admin"
	"OnlineLeadership/internal/usecase/auth"
	"OnlineLeadership/internal/usecase/leaderboard"
	"OnlineLeadership/internal/usecase/score_history"
	"context"
	"github.com/google/uuid"
)

type Auth interface {
	Register(ctx context.Context, user domain.User) (uuid.UUID, error)
	Login(ctx context.Context, username, password string) (string, string, error)
	ParseRefreshToken(ctx context.Context, tokenR string) (string, error)
	ParseAccessToken(ctx context.Context, token string) (uuid.UUID, error)
	GenerateAccessToken(userId string) (string, error)
}
type ScoreHistory interface {
	SubmitScore(ctx context.Context, userID uuid.UUID, gameID uuid.UUID, score int) error
}
type Admin interface {
	Create(ctx context.Context, name string) (uuid.UUID, error)
	GetGames(ctx context.Context) ([]domain.Game, error)
}
type Leaderboard interface {
	GetGlobalLeaderboard(ctx context.Context, offset, limit int) ([]domain.LeaderboardUser, error)
	GetLeaderboard(ctx context.Context, gameID uuid.UUID) ([]domain.LeaderboardUser, error)
	GetMyRank(ctx context.Context, userID uuid.UUID) (int64, error)
}
type Service struct {
	Auth
	ScoreHistory
	Admin
	Leaderboard
}

func NewService(rep *repository.Repository, log *logger.SlogLogger, tokens auth.TokenManager) *Service {
	return &Service{
		Auth:         auth.NewServiceAuth(rep, log, tokens),
		ScoreHistory: score_history.NewScoreService(rep, log),
		Admin:        admin.NewServiceAdmin(rep, log),
		Leaderboard:  leaderboard.NewServiceLeaderboard(rep, log),
	}
}
