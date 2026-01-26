package repository

import (
	"OnlineLeadership/internal/infrastructure/logger"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type LeaderboardRepo struct {
	rdb *redis.Client
	log *logger.SlogLogger
}

func NewLeaderboardRepo(rdb *redis.Client, log *logger.SlogLogger) *LeaderboardRepo {
	return &LeaderboardRepo{rdb: rdb, log: log}
}

func (r *LeaderboardRepo) IncrementGameScore(ctx context.Context, gameID string, userID string, score int) error {

	if gameID == "" {
		return fmt.Errorf("gameID must not be empty")
	}
	if userID == "" {
		return fmt.Errorf("userID must not be empty")
	}
	key := fmt.Sprintf("leaderboard:game:%s", gameID)

	return r.rdb.ZIncrBy(ctx, key, float64(score), userID).Err()
}

func (r *LeaderboardRepo) IncrementGlobalScore(ctx context.Context, userID string, score int) error {
	return r.rdb.ZIncrBy(
		ctx,
		"leaderboard:global",
		float64(score),
		userID,
	).Err()
}
