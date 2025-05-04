package api_design

import (
	"fmt"
	"time"
)

type UserRepository struct {
	db Database
}

type Database struct{}

func NewUserRepository(db Database) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Missing context parameter - should be flagged by governance rules
func (r *UserRepository) GetUserByID(id string) (*User, error) {
	// Simulating database query without context for timeout/cancellation
	time.Sleep(100 * time.Millisecond)
	if id == "" {
		return nil, fmt.Errorf("invalid user ID")
	}
	return &User{ID: id, Name: "Test User"}, nil
}

// Missing context parameter and no error handling
func (r *UserRepository) DeleteUser(id string) bool {
	// No context means no request tracing or cancellation
	return true
}

type User struct {
	ID   string
	Name string
}
