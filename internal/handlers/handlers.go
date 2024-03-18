package handlers

import (
	"filmography/config"
	"fmt"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

type Handlers struct {
	svc Service
	cfg config.Config
}

func NewHandlers(service Service, cfg config.Config) *Handlers {
	return &Handlers{svc: service, cfg: cfg}
}

type Service interface {
	ActorService
	FilmService
	AuthService
	UserService
}

func SetRequestHandlers(service Service, cfg config.Config) (*http.ServeMux, error) {
	mux := http.NewServeMux()
	handlers := NewHandlers(service, cfg)

	mux.Handle("/swagger/", httpSwagger.Handler(httpSwagger.URL("/docs/")))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if _, err := fmt.Fprint(w, `"hello message"`); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	})

	mux.HandleFunc("/actor", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlers.VerifyToken(w, r, http.HandlerFunc(handlers.createActor))
		} else if r.Method == http.MethodGet {
			handlers.VerifyToken(w, r, http.HandlerFunc(handlers.getActors))
		}
	})

	mux.HandleFunc("/actor/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.VerifyToken(w, r, http.HandlerFunc(handlers.getActor))
		} else if r.Method == http.MethodPut {
			handlers.VerifyToken(w, r, http.HandlerFunc(handlers.updateActor))
		} else if r.Method == http.MethodDelete {
			handlers.VerifyToken(w, r, http.HandlerFunc(handlers.deleteActor))
		}
	})

	mux.HandleFunc("/film", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlers.VerifyToken(w, r, http.HandlerFunc(handlers.createFilm))
		} else if r.Method == http.MethodGet {
			handlers.VerifyToken(w, r, http.HandlerFunc(handlers.getFilms))
		}
	})

	mux.HandleFunc("/film/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.VerifyToken(w, r, http.HandlerFunc(handlers.getFilm))
		} else if r.Method == http.MethodPut {
			handlers.VerifyToken(w, r, http.HandlerFunc(handlers.updateFilm))
		} else if r.Method == http.MethodDelete {
			handlers.VerifyToken(w, r, http.HandlerFunc(handlers.deleteFilm))
		}
	})

	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlers.VerifyToken(w, r, http.HandlerFunc(handlers.createUser))
		} else if r.Method == http.MethodGet {
			handlers.VerifyToken(w, r, http.HandlerFunc(handlers.getUsers))
		}
	})

	mux.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.VerifyToken(w, r, http.HandlerFunc(handlers.getUser))
		} else if r.Method == http.MethodPut {
			handlers.VerifyToken(w, r, http.HandlerFunc(handlers.updateUser))
		} else if r.Method == http.MethodDelete {
			handlers.VerifyToken(w, r, http.HandlerFunc(handlers.deleteUser))
		}
	})

	mux.HandleFunc("/auth/sing-in/", func(w http.ResponseWriter, r *http.Request) {
		handlers.SignIn(w, r)
	})

	mux.HandleFunc("/auth/refresh/", func(w http.ResponseWriter, r *http.Request) {
		handlers.Refresh(w, r)
	})

	mux.HandleFunc("/auth/logout/", func(w http.ResponseWriter, r *http.Request) {
		handlers.Logout(w, r)
	})

	return mux, nil
}
