package authentication

import (
	"banners/internal/auth"
	"banners/internal/usecase"
	"context"
	"database/sql"
	"errors"
)

type (
	Authenticator interface {
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
	Authenticator Authenticator
	Repo          Repository
}
type AuthSystem struct {
	Deps
}

func NewAuthenticationSystem(deps Deps) *AuthSystem {
	return &AuthSystem{
		Deps: deps,
	}
}

func (s *AuthSystem) AdminToken(ctx context.Context, credentials usecase.CredentialsDTO) (string, error) {
	expectedPassword, err := s.Repo.AdminCredentials(ctx, credentials.Username)
	if err != nil {
		return "", auth.ErrUnauthorized
	}
	if expectedPassword != credentials.Password {
		return "", auth.ErrForbidden
	}
	token, err := s.Authenticator.AdminToken(ctx, credentials)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *AuthSystem) AdminAuth(ctx context.Context, token string) error {
	return s.Deps.Authenticator.AdminAuth(ctx, token)
}

func (s *AuthSystem) AdminCredentials(ctx context.Context, username string) (string, error) {
	password, err := s.Repo.AdminCredentials(ctx, username)
	if errors.Is(err, sql.ErrNoRows) {
		return "", auth.ErrUnauthorized
	}
	if err != nil {
		return "", err
	}
	return password, nil
}

func (s *AuthSystem) UserToken(ctx context.Context, credentials usecase.CredentialsDTO) (string, error) {
	expectedPassword, err := s.Repo.UserCredentials(ctx, credentials.Username)
	if err != nil {
		return "", auth.ErrUnauthorized
	}
	if expectedPassword != credentials.Password {
		return "", auth.ErrUnauthorized
	}

	token, err := s.Deps.Authenticator.UserToken(ctx, credentials)
	if err != nil {
		return "", err
	}
	return token, nil
}
func (s *AuthSystem) UserAuth(ctx context.Context, token string) error {
	return s.Deps.Authenticator.UserAuth(ctx, token)
}
func (s *AuthSystem) UserCredentials(ctx context.Context, username string) (string, error) {
	password, err := s.Repo.UserCredentials(ctx, username)
	if errors.Is(err, sql.ErrNoRows) {
		return "", auth.ErrUnauthorized
	}
	if err != nil {
		return "", err
	}
	return password, nil
}
