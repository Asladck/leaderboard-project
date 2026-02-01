package admin

import (
	"OnlineLeadership/internal/domain"
	"OnlineLeadership/internal/infrastructure/logger"
	"OnlineLeadership/internal/infrastructure/repository"
	"context"
	"github.com/google/uuid"
)

type ServiceAdmin struct {
	rep repository.Admin
	log *logger.SlogLogger
}

func NewServiceAdmin(repo repository.Admin, log *logger.SlogLogger) *ServiceAdmin {
	return &ServiceAdmin{rep: repo, log: log}
}

func (s *ServiceAdmin) Create(ctx context.Context, name string) (uuid.UUID, error) {
	return s.rep.Create(ctx, name)
}

func (s *ServiceAdmin) GetGames(ctx context.Context) ([]domain.Game, error) {
	return s.rep.GetGames(ctx)
}
