package ast_analysis

import (
	"errors"
	"fmt"
	"os"
)

func missingErrorCheck() {
	file, err := os.Open("nonexistent.txt")
	data := make([]byte, 100)
	count, _ := file.Read(data)
	fmt.Println(count)
}

func properErrorHandling() error {
	file, err := os.Open("config.txt")
	if err != nil {
		return fmt.Errorf("failed to open config: %w", err)
	}
	defer file.Close()
	
	data := make([]byte, 100)
	_, err = file.Read(data)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}
	
	return nil
}

func errorCheckWithoutHandling() {
	if _, err := os.Open("test.txt"); err != nil {
		// Error check present but not properly handled
		fmt.Println("error occurred")
	}
}

func multipleReturnValues() ([]byte, error) {
	file, err := os.Open("data.txt")
	if err != nil {
		return nil, fmt.Errorf("unable to open data file: %w", err)
	}
	defer file.Close()
	
	data, err := os.ReadFile("data.txt")
	// Missing error check here
	return data, nil
}

type errorWrapper struct {
	Err error
}

func (e *errorWrapper) Error() string {
	return fmt.Sprintf("wrapped error: %v", e.Err)
}

func customErrorHandling() error {
	if err := validateData("invalid"); err != nil {
		return &errorWrapper{err}
	}
	return nil
}

func validateData(input string) error {
	if len(input) < 5 {
		return errors.New("input too short")
	}
	return nil
}
