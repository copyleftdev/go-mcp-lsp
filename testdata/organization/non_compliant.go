package organization

import (
	"fmt"
	"math/rand"
)

// Global variables - violates organizational standards
var (
	GlobalCounter int
	GlobalConfig  = make(map[string]string)
)

// Function lacks proper naming convention
func do_something() {
	fmt.Println("Doing something")
}

// Direct dependency without injection
type Service struct{}

func (s *Service) Process() {
	// Direct instantiation of dependencies
	repo := &Repository{}
	repo.Save("data")
}

// Missing interface definition for testability
type Repository struct{}

func (r *Repository) Save(data string) error {
	// Implementation tightly coupled to concrete implementation
	GlobalCounter++
	GlobalConfig["lastSave"] = data
	return nil
}

// Non-deterministic behavior due to random seed
func GenerateID() string {
	return fmt.Sprintf("ID-%d", rand.Intn(1000))
}

// Function with too many parameters
func ComplexOperation(param1, param2, param3, param4, param5, param6, param7, param8 string) {
	fmt.Println("Too many parameters")
}
