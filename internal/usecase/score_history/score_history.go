package score_history

import (
	"OnlineLeadership/internal/infrastructure/logger"
	"context"

	"OnlineLeadership/internal/infrastructure/repository"

	"github.com/google/uuid"
)

type ScoreService struct {
	repo *repository.Repository
	log  *logger.SlogLogger
}

func NewScoreService(repo *repository.Repository, slogLogger *logger.SlogLogger) *ScoreService {
	return &ScoreService{repo: repo, log: slogLogger}
}

func (s *ScoreService) SubmitScore(ctx context.Context, userID uuid.UUID, gameID uuid.UUID, score int) error {
	s.log.Info(ctx, "submit score",
		"user_id", userID,
		"game_id", gameID,
		"score", score,
	)

	// 1️⃣ сохраняем историю (Postgres)
	if err := s.repo.ScoreHistory.Save(ctx, userID, gameID, score); err != nil {
		return err
	}

	// 2️⃣ обновляем leaderboard игры (Redis)
	if err := s.repo.LeaderBoard.IncrementGameScore(ctx, gameID.String(), userID.String(), score); err != nil {
		return err
	}

	// 3️⃣ обновляем глобальный leaderboard
	if err := s.repo.LeaderBoard.IncrementGlobalScore(ctx, userID.String(), score); err != nil {
		return err
	}

	return nil
}
