package ast_analysis

import (
	"context"
	"fmt"
	"time"
)

type UserService struct {
	repo UserRepository
}

type UserRepository interface {
	FindByID(ctx context.Context, id string) (*User, error)
}

type User struct {
	ID   string
	Name string
}

// Good: Has context parameter
func (s *UserService) GetUser(ctx context.Context, id string) (*User, error) {
	return s.repo.FindByID(ctx, id)
}

// Bad: Missing context parameter
func (s *UserService) DeleteUser(id string) error {
	// Should have context as first parameter
	return fmt.Errorf("user %s not found", id)
}

// Good: Function with timeout using context
func (s *UserService) GetUserWithTimeout(ctx context.Context, id string) (*User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	
	return s.repo.FindByID(ctx, id)
}

// Bad: Method with multiple parameters but no context
func (s *UserService) SearchUsers(query string, limit int, offset int) ([]*User, error) {
	// Should have context as first parameter
	return []*User{}, nil
}

type DatabaseClient struct {
	// Some fields
}

// Bad: Missing context in DB operation
func (db *DatabaseClient) ExecuteQuery(query string) error {
	// Should have context as first parameter
	return nil
}

// Good: Proper context usage in DB operation
func (db *DatabaseClient) Execute(ctx context.Context, query string) error {
	// Can respect timeout/cancellation from context
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		// Execute query
		return nil
	}
}
