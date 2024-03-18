package service

import (
	"context"
	"filmography/internal/entities"
	"fmt"
	"github.com/google/uuid"
)

type UserService struct {
	repo UserRepoInterface
}

type UserRepoInterface interface {
	CreateUser(ctx context.Context, user entities.UserEntity) error
	GetUsers(ctx context.Context) ([]entities.UserEntity, error)
	GetUser(ctx context.Context, id string) (entities.UserEntity, error)
	UpdateUser(ctx context.Context, id string, user entities.UserEntity) error
	DeleteUser(ctx context.Context, id string) error
}

func NewUserService(repo UserRepoInterface) UserService {
	return UserService{
		repo: repo,
	}
}

func (svc UserService) CreateUser(ctx context.Context, user entities.UserEntity) error {
	user.ID = uuid.NewString()
	err := svc.repo.CreateUser(ctx, user)
	return err
}

func (svc UserService) GetUsers(ctx context.Context) ([]entities.UserEntity, error) {
	users, err := svc.repo.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("get users failed: %w", err)
	}
	return users, err
}

func (svc UserService) GetUser(ctx context.Context, id string) (entities.UserEntity, error) {
	user, err := svc.repo.GetUser(ctx, id)
	if err != nil {
		return entities.UserEntity{}, fmt.Errorf("get user failed: %w", err)
	}
	return user, err
}

func (svc UserService) UpdateUser(ctx context.Context, id string, user entities.UserEntity) error {
	err := svc.repo.UpdateUser(ctx, id, user)
	return err
}

func (svc UserService) DeleteUser(ctx context.Context, id string) error {
	err := svc.repo.DeleteUser(ctx, id)
	return err
}
