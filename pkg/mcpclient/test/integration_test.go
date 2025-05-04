package test

import (
	"testing"

	"github.com/yourorg/go-mcp-lsp/pkg/mcpclient"
)

func TestClientConnection(t *testing.T) {
	client := mcpclient.New("localhost:9000")
	
	resource, err := client.GetResource("error_handling")
	if err != nil {
		t.Fatalf("Failed to get resource: %v", err)
	}
	
	if resource == nil {
		t.Fatal("Expected resource data, got nil")
	}
}

func TestValidateCode(t *testing.T) {
	client := mcpclient.New("localhost:9000")
	
	code := `package main

import "fmt"

func main() {
	err := doSomething()
	fmt.Println("Done")
}

func doSomething() error {
	return nil
}
`
	
	result, err := client.ValidateCode(code, []string{"error_handling"}, "go")
	if err != nil {
		t.Fatalf("Failed to validate code: %v", err)
	}
	
	if result.Valid {
		t.Error("Expected validation to fail due to missing error handling")
	}
}

func TestGetPrompt(t *testing.T) {
	client := mcpclient.New("localhost:9000")
	
	prompt, err := client.GetPrompt("service implementation", "go", "service")
	if err != nil {
		t.Fatalf("Failed to get prompt: %v", err)
	}
	
	if prompt == "" {
		t.Error("Expected non-empty prompt template")
	}
}

func TestCallTool(t *testing.T) {
	client := mcpclient.New("localhost:9000")
	
	params := map[string]interface{}{
		"template": "go/service",
		"data": map[string]string{
			"PackageName": "users",
			"ServiceName": "UserService",
		},
	}
	
	result, err := client.CallTool("generateScaffold", params, "")
	if err != nil {
		t.Fatalf("Failed to call tool: %v", err)
	}
	
	if result == nil {
		t.Error("Expected result from tool call, got nil")
	}
}
