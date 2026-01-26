package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	accessTTL  = 30 * time.Minute
	refreshTTL = 7 * 24 * time.Hour

	accessTokenType  = "access"
	refreshTokenType = "refresh"

	tokenIssuer = "auth-service"
)

type TokenManager struct {
	accessKey  []byte
	refreshKey []byte
}

func NewTokenManager(accessKey, refreshKey string) *TokenManager {
	return &TokenManager{
		accessKey:  []byte(accessKey),
		refreshKey: []byte(refreshKey),
	}
}

type Claims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
	Type   string `json:"type"`
}

//////////////////////
// TOKEN GENERATION //
//////////////////////

func (m *TokenManager) NewAccessToken(userID string) (string, error) {
	return m.newToken(userID, accessTokenType, accessTTL, m.accessKey)
}

func (m *TokenManager) NewRefreshToken(userID string) (string, error) {
	return m.newToken(userID, refreshTokenType, refreshTTL, m.refreshKey)
}

func (m *TokenManager) newToken(
	userID string,
	tokenType string,
	ttl time.Duration,
	key []byte,
) (string, error) {

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    tokenIssuer,
			Subject:   userID,
		},
		UserID: userID,
		Type:   tokenType,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(key)
}

//////////////////////
// TOKEN PARSING    //
//////////////////////

func (m *TokenManager) ParseAccessToken(tokenStr string) (string, error) {
	return m.parse(tokenStr, accessTokenType, m.accessKey)
}

func (m *TokenManager) ParseRefreshToken(tokenStr string) (string, error) {
	return m.parse(tokenStr, refreshTokenType, m.refreshKey)
}

func (m *TokenManager) parse(
	tokenStr string,
	expectedType string,
	key []byte,
) (string, error) {

	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return key, nil
		},
	)

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}

	if claims.Type != expectedType {
		return "", errors.New("invalid token type")
	}

	if claims.Issuer != tokenIssuer {
		return "", errors.New("invalid token issuer")
	}

	return claims.UserID, nil
}
