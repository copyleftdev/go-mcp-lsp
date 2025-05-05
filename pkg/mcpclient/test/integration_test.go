package test

import (
	"net"
	"os"
	"testing"
	"time"

	"github.com/yourorg/go-mcp-lsp/pkg/mcpclient"
)

// checkServerAvailable tests if the MCP server is running
func checkServerAvailable(address string) bool {
	conn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func TestClientConnection(t *testing.T) {
	serverAddr := "localhost:9000"
	
	// Skip test if CI environment or server not available
	if os.Getenv("CI") != "" || !checkServerAvailable(serverAddr) {
		t.Skip("Skipping integration test: MCP server not available")
	}
	
	client := mcpclient.New(serverAddr)
	
	resource, err := client.GetResource("error_handling")
	if err != nil {
		t.Fatalf("Failed to get resource: %v", err)
	}
	
	if resource == nil {
		t.Fatal("Expected resource data, got nil")
	}
}

func TestValidateCode(t *testing.T) {
	serverAddr := "localhost:9000"
	
	// Skip test if CI environment or server not available
	if os.Getenv("CI") != "" || !checkServerAvailable(serverAddr) {
		t.Skip("Skipping integration test: MCP server not available")
	}
	
	client := mcpclient.New(serverAddr)
	
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
	serverAddr := "localhost:9000"
	
	// Skip test if CI environment or server not available
	if os.Getenv("CI") != "" || !checkServerAvailable(serverAddr) {
		t.Skip("Skipping integration test: MCP server not available")
	}
	
	client := mcpclient.New(serverAddr)
	
	prompt, err := client.GetPrompt("service implementation", "go", "service")
	if err != nil {
		t.Fatalf("Failed to get prompt: %v", err)
	}
	
	if prompt == "" {
		t.Error("Expected non-empty prompt template")
	}
}

func TestCallTool(t *testing.T) {
	serverAddr := "localhost:9000"
	
	// Skip test if CI environment or server not available
	if os.Getenv("CI") != "" || !checkServerAvailable(serverAddr) {
		t.Skip("Skipping integration test: MCP server not available")
	}
	
	client := mcpclient.New(serverAddr)
	
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
