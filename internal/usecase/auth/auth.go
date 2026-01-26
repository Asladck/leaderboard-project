package auth

import (
	"OnlineLeadership/internal/domain"
	"OnlineLeadership/internal/infrastructure/logger"
	"OnlineLeadership/internal/infrastructure/repository"
	"context"
	"golang.org/x/crypto/bcrypt"
)

type TokenManager interface {
	NewAccessToken(userID string) (string, error)
	NewRefreshToken(userID string) (string, error)
	ParseAccessToken(token string) (string, error)
	ParseRefreshToken(token string) (string, error)
}

type Service struct {
	repo   repository.Auth
	log    *logger.SlogLogger
	tokens TokenManager
}

func NewService(repo repository.Auth, log *logger.SlogLogger, tokens TokenManager) *Service {
	return &Service{
		repo:   repo,
		log:    log,
		tokens: tokens,
	}
}

func (s *Service) Register(ctx context.Context, user domain.User) (string, error) {
	hash, err := hashPassword(user.Password)
	if err != nil {
		s.log.Error(ctx, "service auth: hash password error", err.Error())
		return "", err
	}
	user.Password = hash
	return s.repo.CreateUser(ctx, user)
}

func (s *Service) Login(ctx context.Context, username, password string) (string, string, error) {

	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		s.log.Error(ctx, "repo auth: get user error", err.Error())
		return "", "", err
	}

	if err := checkPassword(password, user.Password); err != nil {
		s.log.Error(ctx, "repo auth: check password error", err.Error())

		return "", "", err
	}

	access, err := s.tokens.NewAccessToken(user.Id)
	if err != nil {
		s.log.Error(ctx, "service auth: access error", err.Error())
		return "", "", err
	}

	refresh, err := s.tokens.NewRefreshToken(user.Id)
	if err != nil {
		s.log.Error(ctx, "service auth: ref error", err.Error())
		return "", "", err

	}

	return access, refresh, nil
}

func (s *Service) ParseAccessToken(token string) (string, error) {
	return s.tokens.ParseAccessToken(token)
}

func (s *Service) ParseRefreshToken(token string) (string, error) {
	return s.tokens.ParseRefreshToken(token)
}
func (s *Service) GenerateAccessToken(userId string) (string, error) {
	return s.tokens.NewAccessToken(userId)
}
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
