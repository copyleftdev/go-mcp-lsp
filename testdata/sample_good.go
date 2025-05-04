package testdata

import (
	"fmt"
)

func GoodErrorHandling() error {
	err := riskyOperation()
	if err != nil {
		return fmt.Errorf("operation failed: %w", err)
	}
	return nil
}

func WithContextualError() error {
	if err := complexOperation(); err != nil {
		return fmt.Errorf("complex operation failed: %w", err)
	}
	return nil
}

func complexOperation() error {
	return riskyOperation()
}

func riskyOperation() error {
	return nil
}
