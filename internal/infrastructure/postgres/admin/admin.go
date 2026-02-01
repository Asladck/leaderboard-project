package admin

import (
	"OnlineLeadership/internal/domain"
	"OnlineLeadership/internal/infrastructure/logger"
	"OnlineLeadership/internal/infrastructure/postgres"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RepositoryAdmin struct {
	db  *sqlx.DB
	log *logger.SlogLogger
}

func NewAdminRepository(db *sqlx.DB, log *logger.SlogLogger) *RepositoryAdmin {
	return &RepositoryAdmin{db: db, log: log}
}

func (r *RepositoryAdmin) Create(ctx context.Context, name string) (uuid.UUID, error) {
	var id uuid.UUID
	query := fmt.Sprintf(`INSERT INTO %s (name) VALUES ($1) RETURNING id`, postgres.Games)
	row := r.db.QueryRowContext(ctx, query, name)
	err := row.Scan(&id)
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}
func (r *RepositoryAdmin) GetGames(ctx context.Context) ([]domain.Game, error) {
	var games []domain.Game
	query := fmt.Sprintf(`SELECT * FROM %s`, postgres.Games)
	err := r.db.Select(&games, query)
	if err != nil {
		r.log.Error(ctx, "repository get games error :", err.Error())
		return []domain.Game{}, err
	}
	return games, nil
}
