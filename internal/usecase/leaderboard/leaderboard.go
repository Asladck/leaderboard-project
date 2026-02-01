package leaderboard

import (
	"OnlineLeadership/internal/domain"
	"OnlineLeadership/internal/infrastructure/logger"
	"OnlineLeadership/internal/infrastructure/repository"
	"context"
	"github.com/google/uuid"
)

type ServiceLeaderboard struct {
	repo repository.LeaderBoard
	log  *logger.SlogLogger
}

func NewServiceLeaderboard(repo repository.LeaderBoard, log *logger.SlogLogger) *ServiceLeaderboard {
	return &ServiceLeaderboard{
		repo: repo,
		log:  log,
	}
}

func (s *ServiceLeaderboard) GetGlobalLeaderboard(ctx context.Context, offset int, limit int) ([]domain.LeaderboardUser, error) {
	users, err := s.repo.GetGlobal(ctx, offset, limit)
	if err != nil {
		s.log.Error(ctx, "get global leaderboard error", err.Error())
		return nil, err
	}
	s.log.Info(ctx, "get global service ")
	return users, nil
}
func (s *ServiceLeaderboard) GetMyRank(ctx context.Context, userID uuid.UUID) (int64, error) {
	rank, err := s.repo.GetMyRank(ctx, userID)
	if err != nil {
		s.log.Error(ctx, "repo get my rank error ", err.Error())
		return 0, err
	}
	s.log.Info(ctx, "service get my rank passed")
	return rank, nil
}

func (s *ServiceLeaderboard) GetLeaderboard(ctx context.Context, gameID uuid.UUID) ([]domain.LeaderboardUser, error) {
	users, err := s.repo.GetLeaderboard(ctx, gameID)
	if err != nil {
		s.log.Error(ctx, "repo get leaderboard error", err.Error())
		return []domain.LeaderboardUser{}, err
	}
	s.log.Info(ctx, "service get leaderboard passed")
	return users, nil
}
