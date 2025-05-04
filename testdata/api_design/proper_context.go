package api_design

import (
	"context"
	"fmt"
	"time"
)

type UserService struct {
	repo Repository
}

type UserServiceConfig struct {
	Timeout time.Duration
}

type Repository interface {
	FindByID(ctx context.Context, id string) (*User, error)
	Store(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}

func NewUserService(repo Repository, cfg UserServiceConfig) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) GetUser(ctx context.Context, id string) (*User, error) {
	if id == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}
	
	// Context is properly passed down the call chain
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}
	
	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, user *User) error {
	if user == nil {
		return fmt.Errorf("user cannot be nil")
	}
	
	// Context passes deadlines and cancellation signals
	if err := s.repo.Store(ctx, user); err != nil {
		return fmt.Errorf("failed to store user: %w", err)
	}
	
	return nil
}

func (s *UserService) RemoveUser(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("user ID cannot be empty")
	}
	
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to remove user: %w", err)
	}
	
	return nil
}
