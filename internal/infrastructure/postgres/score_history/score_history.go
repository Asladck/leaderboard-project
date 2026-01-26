package repository

import (
	"OnlineLeadership/internal/infrastructure/logger"
	"context"
	"github.com/jmoiron/sqlx"

	"github.com/google/uuid"
)

type ScoreHistoryRepo struct {
	db  *sqlx.DB
	log *logger.SlogLogger
}

func NewScoreHistoryRepo(db *sqlx.DB, log *logger.SlogLogger) *ScoreHistoryRepo {
	return &ScoreHistoryRepo{db: db, log: log}
}

func (r *ScoreHistoryRepo) Save(ctx context.Context, userID uuid.UUID, gameID uuid.UUID, score int) error {
	query := `
		INSERT INTO score_history (id, user_id, game_id, score)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		uuid.New(),
		userID,
		gameID,
		score,
	)

	return err
}
