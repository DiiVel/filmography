package service

import (
	"filmography/config"
	"filmography/internal/entities"
	"fmt"
	"log"
	"time"
)

type AuthService struct {
	repo TokenRepo
	cfg  config.Config
}

type TokenRepo interface {
	AddToken(token string, expired time.Duration) error
	GetToken(token string) bool
}

func NewAuthService(repo TokenRepo, cfg config.Config) AuthService {
	return AuthService{
		repo: repo,
		cfg:  cfg,
	}
}

func (svc AuthService) SingIn(authInfo entities.Auth) (*entities.Token, error) {
	log.Printf("%+v\n", svc.cfg)
	if authInfo.Login != svc.cfg.AdminLogin || authInfo.Password != svc.cfg.PostgresDBPassword {
		return nil, entities.ErrWrongLoginOrPassword
	}

	params := TokenParams{
		Type:            Admin,
		Hs256Secret:     svc.cfg.Hs256Secret,
		AccessTokenExp:  svc.cfg.AccessTokenExp,
		RefreshTokenExp: svc.cfg.RefreshTokenExp,
	}

	token, err := NewToken(params)
	if err != nil {
		return nil, fmt.Errorf("new token failed: %w", err)
	}

	return token, nil
}

func (svc AuthService) Logout(token string, expired time.Duration) error {
	return svc.repo.AddToken(token, expired)
}

func (svc AuthService) CheckToken(userId string) bool {
	return svc.repo.GetToken(userId)
}
