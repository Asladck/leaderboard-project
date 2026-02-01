package repository

import (
	"OnlineLeadership/internal/domain"
	"OnlineLeadership/internal/infrastructure/logger"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type LeaderboardRepo struct {
	db  *sqlx.DB
	rdb *redis.Client
	log *logger.SlogLogger
}

func NewLeaderboardRepo(db *sqlx.DB, rdb *redis.Client, log *logger.SlogLogger) *LeaderboardRepo {
	return &LeaderboardRepo{db: db, rdb: rdb, log: log}
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
func (r *LeaderboardRepo) GetGlobal(ctx context.Context, offset int, limit int) ([]domain.LeaderboardUser, error) {
	start := int64(offset)
	stop := int64(offset + limit - 1)

	values, err := r.rdb.ZRevRangeWithScores(
		ctx,
		"leaderboard:global",
		start,
		stop,
	).Result()
	if err != nil {
		return nil, err
	}

	result := make([]domain.LeaderboardUser, 0, len(values))

	for _, v := range values {
		userID, err := uuid.Parse(v.Member.(string))
		if err != nil {
			continue
		}

		rank, err := r.rdb.ZRevRank(ctx, "leaderboard:global", v.Member.(string)).Result()
		if err != nil {
			continue
		}

		result = append(result, domain.LeaderboardUser{
			UserID: userID,
			Score:  int64(v.Score),
			Rank:   rank + 1, // ✅ реальный глобальный ранг
		})
	}

	return result, nil
}

func (r *LeaderboardRepo) GetMyRank(ctx context.Context, userID uuid.UUID) (int64, error) {
	rank, err := r.rdb.ZRevRank(ctx, "leaderboard:global", userID.String()).Result()
	if errors.Is(err, redis.Nil) {
		return -1, nil
	}
	if err != nil {
		return 0, err
	}
	return rank + 1, nil
}
func (r *LeaderboardRepo) GetLeaderboard(ctx context.Context, gameID uuid.UUID) ([]domain.LeaderboardUser, error) {
	if gameID == uuid.Nil {
		return nil, fmt.Errorf("gameID must not be empty")
	}
	key := fmt.Sprintf("leaderboard:game:%s", gameID.String())

	values, err := r.rdb.ZRevRangeWithScores(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	result := make([]domain.LeaderboardUser, 0, len(values))
	for _, v := range values {
		memberStr, ok := v.Member.(string)
		if !ok {
			continue
		}

		userID, err := uuid.Parse(memberStr)
		if err != nil {
			continue
		}

		rank, err := r.rdb.ZRevRank(ctx, key, memberStr).Result()
		if err != nil {
			continue
		}

		result = append(result, domain.LeaderboardUser{
			UserID: userID,
			Score:  int64(v.Score),
			Rank:   rank + 1,
		})
	}

	return result, nil
}
