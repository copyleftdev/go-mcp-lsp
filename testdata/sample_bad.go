package testdata

import "fmt"

func BadErrorHandling() error {
	err := riskyOperation()
	return nil
}

func IgnoredErrors() {
	_ = riskyOperation()
	fmt.Println("Continuing regardless of errors")
}

func riskyOperation() error {
	return fmt.Errorf("something went wrong")
}
