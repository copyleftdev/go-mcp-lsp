package error_handling

import (
	"fmt"
	"io/ioutil"
	"os"
)

// PoorErrorHandling demonstrates an anti-pattern where errors are not checked
func PoorErrorHandling() {
	data, err := ioutil.ReadFile("nonexistent.txt")
	fmt.Println(string(data))
}

// ImproperPropagation demonstrates incorrect error propagation
func ImproperPropagation() error {
	err := writeToFile("test.txt", "content")
	if err != nil {
		fmt.Println("Failed to write to file:", err)
		// Missing return statement with error
	}
	return nil
}

// IgnoredReturn ignores the error returned by a function
func IgnoredReturn() {
	_ = writeToFile("test.txt", "content")
}

func writeToFile(filename, content string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}
