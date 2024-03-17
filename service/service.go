package service

type Service struct {
	ActorService
	FilmService
}

type Repo interface {
	ActorRepoInterface
	FilmRepoInterface
}

func New(repo Repo) Service {
	return Service{
		ActorService: NewActorService(repo),
		FilmService:  NewFilmService(repo),
	}
}
