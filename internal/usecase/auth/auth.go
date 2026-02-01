package auth

import (
	"OnlineLeadership/internal/domain"
	"OnlineLeadership/internal/infrastructure/logger"
	"OnlineLeadership/internal/infrastructure/repository"
	"context"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TokenManager interface {
	NewAccessToken(userID string) (string, error)
	NewRefreshToken(userID string) (string, error)
	ParseAccessToken(ctx context.Context, token string) (string, error)
	ParseRefreshToken(ctx context.Context, token string) (string, error)
}

type ServiceAuth struct {
	repo   repository.Auth
	log    *logger.SlogLogger
	tokens TokenManager
}

func NewServiceAuth(repo repository.Auth, log *logger.SlogLogger, tokens TokenManager) *ServiceAuth {
	return &ServiceAuth{
		repo:   repo,
		log:    log,
		tokens: tokens,
	}
}

func (s *ServiceAuth) Register(ctx context.Context, user domain.User) (uuid.UUID, error) {
	hash, err := hashPassword(user.Password)
	if err != nil {
		s.log.Error(ctx, "service auth: hash password error", err.Error())
		return uuid.UUID{}, err
	}
	user.Password = hash
	return s.repo.CreateUser(ctx, user)
}

func (s *ServiceAuth) Login(ctx context.Context, username, password string) (string, string, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		s.log.Error(ctx, "repo auth: get user error", err.Error())
		return "", "", err
	}

	if err := checkPassword(password, user.Password); err != nil {
		s.log.Error(ctx, "repo auth: check password error", err.Error())

		return "", "", err
	}

	// Convert uuid.UUID to string for JWT token
	access, err := s.tokens.NewAccessToken(user.Id.String())
	if err != nil {
		s.log.Error(ctx, "service auth: access error", err.Error())
		return "", "", err
	}

	refresh, err := s.tokens.NewRefreshToken(user.Id.String())
	if err != nil {
		s.log.Error(ctx, "service auth: ref error", err.Error())
		return "", "", err

	}

	return access, refresh, nil
}

func (s *ServiceAuth) ParseAccessToken(ctx context.Context, token string) (uuid.UUID, error) {
	userIdStr, err := s.tokens.ParseAccessToken(ctx, token)
	if err != nil {
		s.log.Error(ctx, "parse token error", err.Error())
		return uuid.UUID{}, err
	}

	userID, err := uuid.Parse(userIdStr)
	if err != nil {
		return uuid.UUID{}, err
	}

	return userID, nil
}

func (s *ServiceAuth) ParseRefreshToken(ctx context.Context, token string) (string, error) {
	return s.tokens.ParseRefreshToken(ctx, token)
}
func (s *ServiceAuth) GenerateAccessToken(userId string) (string, error) {
	return s.tokens.NewAccessToken(userId)
}
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
