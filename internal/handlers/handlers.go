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
}

func SetRequestHandlers(service Service, cfg config.Config) (*http.ServeMux, error) {
	mux := http.NewServeMux()
	handlers := NewHandlers(service, cfg)

	mux.Handle("/swagger/", httpSwagger.Handler(httpSwagger.URL("/docs/swagger.yaml")))

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
			handlers.createActor(w, r)
		} else if r.Method == http.MethodGet {
			handlers.getActors(w, r)
		}
	})

	mux.HandleFunc("/actor/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.getActor(w, r)
		} else if r.Method == http.MethodPut {
			handlers.updateActor(w, r)
		} else if r.Method == http.MethodDelete {
			handlers.deleteActor(w, r)
		}
	})

	mux.HandleFunc("/film", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlers.createFilm(w, r)
		} else if r.Method == http.MethodGet {
			handlers.getFilms(w, r)
		}
	})

	mux.HandleFunc("/film/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.getFilm(w, r)
		} else if r.Method == http.MethodPut {
			handlers.updateFilm(w, r)
		} else if r.Method == http.MethodDelete {
			handlers.deleteFilm(w, r)
		}
	})

	return mux, nil
}
