package error_handling

import (
	"errors"
	"fmt"
	"os"
)

func ProperErrorHandling() ([]byte, error) {
	data, err := os.ReadFile("config.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	return data, nil
}

func WithErrorWrapping() error {
	err := writeToConfigFile("settings.txt", "value=true")
	if err != nil {
		return fmt.Errorf("config update failed: %w", err)
	}
	return nil
}

func WithCustomErrors() error {
	if err := validateInput("test"); err != nil {
		return NewValidationError("input validation failed", err)
	}
	return nil
}

type ValidationError struct {
	Msg string
	Err error
}

func (e *ValidationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Msg, e.Err)
	}
	return e.Msg
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}

func NewValidationError(msg string, err error) error {
	return &ValidationError{
		Msg: msg,
		Err: err,
	}
}

func writeToConfigFile(filename, content string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}

func validateInput(input string) error {
	if len(input) < 5 {
		return errors.New("input too short")
	}
	return nil
}
