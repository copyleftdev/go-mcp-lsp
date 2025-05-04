package ast_analysis

import (
	"crypto/md5"  // Intentionally using weak crypto
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"net/http"
	"sync"
)

// Bad: Unsynchronized map access from goroutines
type UnsafeCache struct {
	data map[string]interface{}
}

func NewUnsafeCache() *UnsafeCache {
	return &UnsafeCache{
		data: make(map[string]interface{}),
	}
}

func (c *UnsafeCache) LoadConcurrently(keys []string) {
	for _, key := range keys {
		go func(k string) {
			c.data[k] = "value" // Unsafe concurrent map access
		}(key)
	}
}

// Good: Synchronized map access
type SafeCache struct {
	mu   sync.RWMutex
	data map[string]interface{}
}

func NewSafeCache() *SafeCache {
	return &SafeCache{
		data: make(map[string]interface{}),
	}
}

func (c *SafeCache) LoadConcurrently(keys []string) {
	for _, key := range keys {
		go func(k string) {
			c.mu.Lock()
			c.data[k] = "value"
			c.mu.Unlock()
		}(key)
	}
}

// Bad: Using weak cryptography
func generateWeakHash(data []byte) string {
	hash := md5.Sum(data) // MD5 is cryptographically broken
	return fmt.Sprintf("%x", hash)
}

// Good: Using strong cryptography
func generateStrongHash(data []byte) string {
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}

// Bad: SQL injection vulnerability
func buildUnsafeQuery(username string) string {
	return fmt.Sprintf("SELECT * FROM users WHERE username = '%s'", username)
}

// Good: Parameterized query
func buildSafeQuery() string {
	return "SELECT * FROM users WHERE username = ?"
}

// Bad: Hardcoded credentials
func connectToDatabase() {
	username := "admin"
	password := "super_secret_password" // Hardcoded credential
	
	fmt.Printf("Connecting with %s:%s\n", username, password)
}

// Good: Environment-based credentials
func safeConnectToDatabase() error {
	username := getEnvVar("DB_USER")
	password := getEnvVar("DB_PASSWORD")
	
	if username == "" || password == "" {
		return fmt.Errorf("missing database credentials")
	}
	
	fmt.Println("Connecting to database...")
	return nil
}

func getEnvVar(name string) string {
	// Simplified for testing purposes
	return ""
}

// Bad: XSS vulnerability
func handleUserInput(w http.ResponseWriter, r *http.Request) {
	input := r.URL.Query().Get("name")
	fmt.Fprintf(w, "<div>Hello, %s!</div>", input) // Direct interpolation
}

// Good: XSS prevention
func safeHandleUserInput(w http.ResponseWriter, r *http.Request) {
	input := r.URL.Query().Get("name")
	// In a real implementation, we would use proper HTML escaping
	// or a templating library that handles this automatically
	safeInput := sanitizeHTML(input)
	fmt.Fprintf(w, "<div>Hello, %s!</div>", safeInput)
}

func sanitizeHTML(input string) string {
	// Simple sanitization for demonstration
	return input // In reality, this would do actual sanitization
}
