package service

import (
	"filmography/internal/entities"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	Admin = "admin"
)

var (
	ErrTokenExpired = fmt.Errorf("token expired")
	ErrUnknownType  = fmt.Errorf("unknown type")
)

type TokenParams struct {
	ID              any
	Type            string
	Hs256Secret     string
	AccessTokenExp  int
	RefreshTokenExp int
}

func NewToken(params TokenParams) (*entities.Token, error) {
	if params.Type != Admin {
		return nil, ErrUnknownType
	}

	accessExp := time.Now().Add(time.Duration(params.AccessTokenExp) * time.Minute)

	access, err := newJwt(accessExp, params)
	if err != nil {
		return nil, fmt.Errorf("new jwt failed: %w", err)
	}

	rtExp := time.Now().Add(time.Duration(params.RefreshTokenExp) * 24 * time.Hour)

	rt, err := newJwt(rtExp, params)
	if err != nil {
		return nil, fmt.Errorf("new rt failed: %w", err)
	}

	return &entities.Token{Access: access, RT: rt}, nil
}

func newJwt(jwtExp time.Time, p TokenParams) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["type"] = p.Type
	claims["exp"] = jwtExp.UTC().Unix()

	secret := []byte(p.Hs256Secret)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("signed string failed: %w", err)
	}

	return tokenString, nil
}

func (svc AuthService) Verify(token string) error {
	tokenJwt, err := jwt.Parse(
		token,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(svc.cfg.Hs256Secret), nil
		},
	)

	if err != nil {
		return fmt.Errorf("token parse failed: %w", err)
	}

	claims, ok := tokenJwt.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("jwt map claims failed")
	}

	if !claims.VerifyExpiresAt(time.Now().UTC().Unix(), true) {
		return ErrTokenExpired
	}
	if string(claims["type"].(string)) != Admin {
		return ErrUnknownType
	}
	return nil
}
