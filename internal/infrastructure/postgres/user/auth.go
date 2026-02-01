package user

import (
	"OnlineLeadership/internal/domain"
	"OnlineLeadership/internal/infrastructure/logger"
	"OnlineLeadership/internal/infrastructure/postgres"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Auth struct {
	db  *sqlx.DB
	log *logger.SlogLogger
}

func NewAuthRepository(db *sqlx.DB, log *logger.SlogLogger) *Auth {
	return &Auth{
		db:  db,
		log: log,
	}
}

func (r *Auth) CreateUser(ctx context.Context, user domain.User) (uuid.UUID, error) {
	var id uuid.UUID
	query := fmt.Sprintf("INSERT INTO %s (username,password_hash,email) values ($1,$2,$3) RETURNING id", postgres.Users)
	row := r.db.QueryRow(query, user.Username, user.Password, user.Email)
	if err := row.Scan(&id); err != nil {
		r.log.Error(ctx, "creating user error ", err.Error())
		return uuid.UUID{}, err
	}
	r.log.Info(ctx, "creating user successfully ")
	return id, nil
}

func (r *Auth) GetUser(ctx context.Context, username, password string) (domain.User, error) {
	var user domain.User
	r.log.Info(ctx, username, password)

	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", postgres.Users)

	err := r.db.Get(&user, query, username, password)
	return user, err
}
func (r *Auth) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	var user domain.User
	query := fmt.Sprintf("SELECT id, username, password_hash FROM %s WHERE username=$1", postgres.Users)
	err := r.db.QueryRowContext(ctx, query, username).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		r.log.Error(ctx, "postgres error", err.Error())
	}
	return user, err
}
