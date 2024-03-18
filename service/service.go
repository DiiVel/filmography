package service

import (
	"filmography/config"
)

type Service struct {
	ActorService
	FilmService
	AuthService
	UserService
}

type Repo interface {
	ActorRepoInterface
	FilmRepoInterface
	UserRepoInterface
}

type Cache interface {
	TokenRepo
}

func New(repo Repo, cache Cache, cfg config.Config) Service {
	return Service{
		ActorService: NewActorService(repo),
		FilmService:  NewFilmService(repo),
		AuthService:  NewAuthService(cache, cfg),
	}
}
