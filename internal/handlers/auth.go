package handlers

import (
	"encoding/json"
	"errors"
	"filmography/internal/entities"
	"filmography/service"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"strings"
	"time"
)

type AuthService interface {
	SingIn(authInfo entities.Auth) (*entities.Token, error)
	Verify(token string) error
	Logout(token string, expired time.Duration) error
	CheckToken(userId string) bool
}

// @Summary User sign-in
// @Description Allows a user to sign in with their credentials.
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body entities.Auth true "User credentials"
// @Success 200 {object} map[string]string "Access token response"
// @Failure 400 {string} string "Bad request"
// @Failure 403 {string} string "Forbidden"
// @Failure 500 {string} string "Internal server error"
// @Router /admin/auth/signin [post]
func (handlers Handlers) SignIn(w http.ResponseWriter, r *http.Request) {
	var auth entities.Auth

	err := json.NewDecoder(r.Body).Decode(&auth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println(auth)

	token, err := handlers.svc.SingIn(auth)
	if err != nil {
		if errors.Is(err, entities.ErrWrongLoginOrPassword) {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("sing in failed")
		return
	}

	exp := time.Now().Add(time.Duration(handlers.cfg.RefreshTokenExp) * 24 * time.Hour)
	cookie := http.Cookie{
		Name:    "refresh_token",
		Value:   token.RT,
		Path:    "/admin/auth",
		Expires: exp,
	}
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"access_token": token.Access,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}

func (handlers Handlers) VerifyToken(w http.ResponseWriter, r *http.Request, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.Split(r.Header.Get("Authorization"), " ")
		if len(token) < 2 {
			http.Error(w, fmt.Errorf("access token required").Error(), http.StatusUnauthorized)
			return
		}
		accessToken := token[1]

		err := handlers.svc.Verify(accessToken)
		if err != nil {
			if errors.Is(err, service.ErrTokenExpired) {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			if strings.Contains(err.Error(), jwt.ErrSignatureInvalid.Error()) {
				http.Error(w, fmt.Errorf("wrong signature").Error(), http.StatusForbidden)
				return
			}

			http.Error(w, fmt.Errorf("wrong token").Error(), http.StatusForbidden)
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("verify failed")
			return
		}

		if !handlers.svc.CheckToken(accessToken) {
			http.Error(w, "you already logged out", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// @Summary Refresh access token
// @Description Allows refreshing an access token using a valid refresh token.
// @Tags Auth
// @Produce json
// @Success 200 {object} map[string]string "New access token response"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 500 {string} string "Internal server error"
// @Router /admin/auth/refresh [get]
func (handlers Handlers) Refresh(w http.ResponseWriter, r *http.Request) {
	refresh, err := r.Cookie("refresh_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			http.Error(w, fmt.Errorf("bad refresh token").Error(), http.StatusForbidden)
			return
		}

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("get cookies failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = handlers.svc.Verify(refresh.Value)
	if err != nil {
		if errors.Is(err, service.ErrTokenExpired) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if strings.Contains(err.Error(), jwt.ErrSignatureInvalid.Error()) {
			http.Error(w, fmt.Errorf("wrong signature").Error(), http.StatusForbidden)
			return
		}

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("verify failed")
		http.Error(w, fmt.Errorf("verify rt failed: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	params := service.TokenParams{
		Type:            service.Admin,
		Hs256Secret:     handlers.cfg.Hs256Secret,
		AccessTokenExp:  handlers.cfg.AccessTokenExp,
		RefreshTokenExp: handlers.cfg.RefreshTokenExp,
	}

	token, err := service.NewToken(params)
	if err != nil {
		if errors.Is(err, service.ErrUnknownType) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("new token failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	exp := time.Now().Add(time.Duration(handlers.cfg.RefreshTokenExp) * 24 * time.Hour)
	cookie := http.Cookie{
		Name:  "refresh_token",
		Value: token.RT,
		Path:  "/admin/auth",
		// Set other cookie properties like HttpOnly and Secure based on your needs
		Expires: exp,
	}
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"access_token": token.Access,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}

// @Summary User logout
// @Description Allows a user to log out and invalidate their access token.
// @Tags Auth
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {string} string "Logout successful"
// @Failure 500 {string} string "Internal server error"
// @Router /admin/auth/logout [post]
func (handlers Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	exp := time.Duration(handlers.cfg.AccessTokenExp) * time.Minute
	token := strings.Split(r.Header.Get("Authorization"), " ")
	if len(token) < 2 {
		http.Error(w, fmt.Errorf("access token required").Error(), http.StatusUnauthorized)
		return
	}
	accessToken := token[1]

	err := handlers.svc.Logout(accessToken, exp)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("logout failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:    "refresh_token",
		Value:   "",
		Path:    "/admin/auth",
		Expires: time.Now().Add(-1 * time.Second), // Expire the cookie immediately
	}
	http.SetCookie(w, &cookie)

	w.WriteHeader(http.StatusOK)
}
