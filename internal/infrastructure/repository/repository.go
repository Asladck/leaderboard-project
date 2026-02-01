package repository

import (
	"OnlineLeadership/internal/domain"
	"OnlineLeadership/internal/infrastructure/logger"
	"OnlineLeadership/internal/infrastructure/postgres/admin"
	leader "OnlineLeadership/internal/infrastructure/postgres/leaderboard"
	score "OnlineLeadership/internal/infrastructure/postgres/score_history"
	"OnlineLeadership/internal/infrastructure/postgres/user"
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Auth interface {
	CreateUser(ctx context.Context, user domain.User) (uuid.UUID, error)
	GetUser(ctx context.Context, username, password string) (domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (domain.User, error)
}
type ScoreHistory interface {
	Save(ctx context.Context, userID uuid.UUID, gameID uuid.UUID, score int) error
}
type LeaderBoard interface {
	IncrementGameScore(ctx context.Context, gameID string, userID string, score int) error
	IncrementGlobalScore(ctx context.Context, userID string, score int) error
	GetGlobal(ctx context.Context, offset int, limit int) ([]domain.LeaderboardUser, error)
	GetMyRank(ctx context.Context, userID uuid.UUID) (int64, error)
	GetLeaderboard(ctx context.Context, gameID uuid.UUID) ([]domain.LeaderboardUser, error)
}
type Admin interface {
	Create(ctx context.Context, name string) (uuid.UUID, error)
	GetGames(ctx context.Context) ([]domain.Game, error)
}
type Repository struct {
	Auth
	ScoreHistory
	LeaderBoard
	Admin
}

func NewRepository(db *sqlx.DB, redis *redis.Client, log *logger.SlogLogger) *Repository {
	return &Repository{
		Auth:         user.NewAuthRepository(db, log),
		ScoreHistory: score.NewScoreHistoryRepo(db, log),
		LeaderBoard:  leader.NewLeaderboardRepo(db, redis, log),
		Admin:        admin.NewAdminRepository(db, log),
	}

}
