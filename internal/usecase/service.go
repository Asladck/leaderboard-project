package usecase

import (
	"OnlineLeadership/internal/domain"
	"OnlineLeadership/internal/infrastructure/logger"
	"OnlineLeadership/internal/infrastructure/repository"
	"OnlineLeadership/internal/usecase/auth"
	"OnlineLeadership/internal/usecase/score_history"
	"context"
	"github.com/google/uuid"
)

type Auth interface {
	Register(ctx context.Context, user domain.User) (string, error)
	Login(ctx context.Context, username, password string) (string, string, error)
	ParseRefreshToken(tokenR string) (string, error)
	ParseAccessToken(token string) (string, error)
	GenerateAccessToken(userId string) (string, error)
}
type ScoreHistory interface {
	SubmitScore(ctx context.Context, userID uuid.UUID, gameID uuid.UUID, score int) error
}
type Service struct {
	Auth
	ScoreHistory
}

func NewService(rep *repository.Repository, log *logger.SlogLogger, tokens auth.TokenManager) *Service {
	return &Service{
		Auth:         auth.NewService(rep, log, tokens),
		ScoreHistory: score_history.NewScoreService(rep, log),
	}
}
