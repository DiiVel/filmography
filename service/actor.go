package service

import (
	"context"
	"filmography/internal/entities"
	"fmt"
	"github.com/google/uuid"
)

type ActorService struct {
	repo ActorRepoInterface
}

type ActorRepoInterface interface {
	CreateActor(ctx context.Context, actor entities.ActorEntity) error
	GetActors(ctx context.Context) ([]entities.ActorEntity, error)
	GetActor(ctx context.Context, id string) (entities.ActorEntity, error)
	UpdateActor(ctx context.Context, id string, actor entities.ActorEntity) error
	DeleteActor(ctx context.Context, id string) error
}

func NewActorService(repo ActorRepoInterface) ActorService {
	return ActorService{
		repo: repo,
	}
}

func (svc ActorService) CreateActor(ctx context.Context, actor entities.ActorEntity) error {
	actor.ID = uuid.NewString()
	err := svc.repo.CreateActor(ctx, actor)
	return err
}

func (svc ActorService) GetActors(ctx context.Context) ([]entities.ActorEntity, error) {
	actors, err := svc.repo.GetActors(ctx)
	if err != nil {
		return nil, fmt.Errorf("get actors failed: %w", err)
	}
	return actors, err
}

func (svc ActorService) GetActor(ctx context.Context, id string) (entities.ActorEntity, error) {
	actor, err := svc.repo.GetActor(ctx, id)
	if err != nil {
		return entities.ActorEntity{}, fmt.Errorf("get actor failed: %w", err)
	}
	return actor, err
}

func (svc ActorService) UpdateActor(ctx context.Context, id string, actor entities.ActorEntity) error {
	err := svc.repo.UpdateActor(ctx, id, actor)
	return err
}

func (svc ActorService) DeleteActor(ctx context.Context, id string) error {
	err := svc.repo.DeleteActor(ctx, id)
	return err
}
