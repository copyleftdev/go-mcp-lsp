package organization

import (
	"crypto/rand"
	"fmt"
	"time"
)

type Config struct {
	Settings map[string]string
}

func DoSomething() {
	fmt.Println("Doing something")
}

type RepositoryInterface interface {
	Save(data string) error
	Load(id string) (string, error)
}

type ServiceConfig struct {
	Repository RepositoryInterface
	Timeout    time.Duration
}

type Service struct {
	repo    RepositoryInterface
	timeout time.Duration
}

func NewService(cfg ServiceConfig) *Service {
	return &Service{
		repo:    cfg.Repository,
		timeout: cfg.Timeout,
	}
}

func (s *Service) Process(data string) error {
	return s.repo.Save(data)
}

type Repository struct {
	counter int
	config  map[string]string
}

func NewRepository() *Repository {
	return &Repository{
		config: make(map[string]string),
	}
}

func (r *Repository) Save(data string) error {
	r.counter++
	r.config["lastSave"] = data
	return nil
}

func (r *Repository) Load(id string) (string, error) {
	data, ok := r.config[id]
	if !ok {
		return "", fmt.Errorf("data not found for id: %s", id)
	}
	return data, nil
}

func GenerateID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate ID: %w", err)
	}
	return fmt.Sprintf("ID-%x", b), nil
}

type OperationParams struct {
	UserID      string
	Action      string
	Permissions []string
	Timestamp   time.Time
	Metadata    map[string]string
}

func ComplexOperation(params OperationParams) error {
	if params.UserID == "" {
		return fmt.Errorf("user ID is required")
	}
	return nil
}
