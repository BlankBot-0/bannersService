package authentification

import (
	"banners/internal/config"
	"banners/internal/usecase"
	"context"
	"database/sql"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type (
	Authentificator interface {
		UserToken(ctx context.Context, credentials usecase.CredentialsDTO) (string, error)
		UserAuth(ctx context.Context, token string) error
		AdminToken(ctx context.Context, credentials usecase.CredentialsDTO) (string, error)
		AdminAuth(ctx context.Context, token string) error
	}
	Repository interface {
		AdminCredentials(ctx context.Context, username string) (string, error)
		UserCredentials(ctx context.Context, username string) (string, error)
	}
)

type Deps struct {
	authentificator Authentificator
	Repo            Repository
}
type AuthentificationSystem struct {
	Deps
	PrivateKey     string
	ExpirationTime time.Duration
}

func NewAuthentificationSystem(deps Deps, cfg config.Auth) *AuthentificationSystem {
	return &AuthentificationSystem{
		Deps:           deps,
		PrivateKey:     cfg.PrivateKey,
		ExpirationTime: cfg.ExpirationTime,
	}
}

func (s *AuthentificationSystem) AdminToken(ctx context.Context, credentials usecase.CredentialsDTO) (string, error) {
	expectedPassword, err := s.Repo.AdminCredentials(ctx, credentials.Username)
	if err != nil {
		return "", ErrUnauthorized
	}
	if expectedPassword != credentials.Password {
		return "", ErrForbidden
	}
	expirationTime := time.Now().Add(s.ExpirationTime)
	claims := &AdminClaims{
		Admin: true,
		Claims: Claims{
			Username: credentials.Username,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.PrivateKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *AuthentificationSystem) AdminAuth(ctx context.Context, token string) error {
	claims := &AdminClaims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return s.PrivateKey, nil
	})
	if errors.Is(err, jwt.ErrSignatureInvalid) {
		return ErrForbidden
	}
	if err != nil {
		return ErrInvalidToken
	}
	if !tkn.Valid {
		return ErrUnauthorized
	}
	return nil
}

func (s *AuthentificationSystem) AdminCredentials(ctx context.Context, username string) (string, error) {
	password, err := s.Repo.AdminCredentials(ctx, username)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrUnauthorized
	}
	if err != nil {
		return "", err
	}
	return password, nil
}

func (s *AuthentificationSystem) UserToken(ctx context.Context, credentials usecase.CredentialsDTO) (string, error) {
	expectedPassword, err := s.Repo.UserCredentials(ctx, credentials.Username)
	if err != nil {
		return "", ErrUnauthorized
	}
	if expectedPassword != credentials.Password {
		return "", ErrUnauthorized
	}
	expirationTime := time.Now().Add(s.ExpirationTime)
	claims := &Claims{
		Username: credentials.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.PrivateKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func (s *AuthentificationSystem) UserAuth(ctx context.Context, token string) error {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return s.PrivateKey, nil
	})
	if errors.Is(err, jwt.ErrSignatureInvalid) {
		return ErrUnauthorized
	}
	if err != nil {
		return ErrInvalidToken
	}
	if !tkn.Valid {
		return ErrUnauthorized
	}
	return nil
}
func (s *AuthentificationSystem) UserCredentials(ctx context.Context, username string) (string, error) {
	password, err := s.Repo.UserCredentials(ctx, username)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrUnauthorized
	}
	if err != nil {
		return "", err
	}
	return password, nil
}

var admins = map[string]string{}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type AdminClaims struct {
	Admin bool `json:"admin"`
	Claims
}
