package ast

import (
	"testing"
)

func TestAnalyzeErrorHandling(t *testing.T) {
	tests := []struct {
		name          string
		code          string
		expectedIssue bool
	}{
		{
			name: "Missing error check",
			code: `package test
func foo() {
	file, err := os.Open("test.txt")
	data := make([]byte, 100)
	file.Read(data)
}`,
			expectedIssue: true,
		},
		{
			name: "Proper error handling",
			code: `package test
func foo() error {
	file, err := os.Open("test.txt")
	if err != nil {
		return err
	}
	data := make([]byte, 100)
	_, err = file.Read(data)
	if err != nil {
		return err
	}
	return nil
}`,
			expectedIssue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzer := NewAnalyzer(AnalyzerConfig{IncludeTests: true})
			file, err := analyzer.ParseString("test.go", tt.code)
			if err != nil {
				t.Fatalf("Failed to parse code: %v", err)
			}

			issues := analyzer.AnalyzeErrorHandling(file)
			hasIssue := len(issues) > 0

			if hasIssue != tt.expectedIssue {
				t.Errorf("Expected issue: %v, got: %v", tt.expectedIssue, hasIssue)
				if hasIssue {
					for i, issue := range issues {
						t.Logf("Issue %d: %s at line %d", i+1, issue.Description, issue.Position.Line)
					}
				}
			}
		})
	}
}

func TestAnalyzeAPIDesign(t *testing.T) {
	tests := []struct {
		name          string
		code          string
		expectedIssue bool
	}{
		{
			name: "Missing context parameter",
			code: `package test
type Service struct{}

func (s *Service) DoSomething(id string) error {
	return nil
}`,
			expectedIssue: true,
		},
		{
			name: "With context parameter",
			code: `package test
import "context"

type Service struct{}

func (s *Service) DoSomething(ctx context.Context, id string) error {
	return nil
}`,
			expectedIssue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzer := NewAnalyzer(AnalyzerConfig{IncludeTests: true})
			file, err := analyzer.ParseString("test.go", tt.code)
			if err != nil {
				t.Fatalf("Failed to parse code: %v", err)
			}

			issues := analyzer.AnalyzeAPIDesign(file)
			hasIssue := len(issues) > 0

			if hasIssue != tt.expectedIssue {
				t.Errorf("Expected issue: %v, got: %v", tt.expectedIssue, hasIssue)
				if hasIssue {
					for i, issue := range issues {
						t.Logf("Issue %d: %s at line %d", i+1, issue.Description, issue.Position.Line)
					}
				}
			}
		})
	}
}

func TestAnalyzeConcurrencySafety(t *testing.T) {
	tests := []struct {
		name          string
		code          string
		expectedIssue bool
	}{
		{
			name: "Unsafe concurrent map access",
			code: `package test
func unsafeAccess() {
	m := make(map[string]string)
	go func() {
		m["key"] = "value"
	}()
}`,
			expectedIssue: true,
		},
		{
			name: "Safe concurrent access with mutex",
			code: `package test
import "sync"

func safeAccess() {
	var mu sync.Mutex
	m := make(map[string]string)
	
	go func() {
		mu.Lock()
		m["key"] = "value"
		mu.Unlock()
	}()
}`,
			expectedIssue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzer := NewAnalyzer(AnalyzerConfig{IncludeTests: true})
			file, err := analyzer.ParseString("test.go", tt.code)
			if err != nil {
				t.Fatalf("Failed to parse code: %v", err)
			}

			issues := analyzer.AnalyzeConcurrencySafety(file)
			hasIssue := len(issues) > 0

			if hasIssue != tt.expectedIssue {
				t.Errorf("Expected issue: %v, got: %v", tt.expectedIssue, hasIssue)
				if hasIssue {
					for i, issue := range issues {
						t.Logf("Issue %d: %s at line %d", i+1, issue.Description, issue.Position.Line)
					}
				}
			}
		})
	}
}

func TestAnalyzeSecurityIssues(t *testing.T) {
	tests := []struct {
		name          string
		code          string
		expectedIssue bool
	}{
		{
			name: "Weak cryptography",
			code: `package test
import "crypto/md5"

func generateHash(data []byte) []byte {
	hash := md5.Sum(data)
	return hash[:]
}`,
			expectedIssue: true,
		},
		{
			name: "SQL injection vulnerability",
			code: `package test
import "fmt"
import "database/sql"

func queryUser(username string, db *sql.DB) {
	query := fmt.Sprintf("SELECT * FROM users WHERE username='%s'", username)
	db.Query(query)
}`,
			expectedIssue: true,
		},
		{
			name: "Secure practice",
			code: `package test
import "crypto/sha256"
import "database/sql"

func generateHash(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func queryUser(username string, db *sql.DB) {
	db.Query("SELECT * FROM users WHERE username=?", username)
}`,
			expectedIssue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzer := NewAnalyzer(AnalyzerConfig{IncludeTests: true})
			file, err := analyzer.ParseString("test.go", tt.code)
			if err != nil {
				t.Fatalf("Failed to parse code: %v", err)
			}

			issues := analyzer.AnalyzeSecurityIssues(file)
			hasIssue := len(issues) > 0

			if hasIssue != tt.expectedIssue {
				t.Errorf("Expected issue: %v, got: %v", tt.expectedIssue, hasIssue)
				if hasIssue {
					for i, issue := range issues {
						t.Logf("Issue %d: %s at line %d", i+1, issue.Description, issue.Position.Line)
					}
				}
			}
		})
	}
}

func TestAnalyzeOrganizationStandards(t *testing.T) {
	tests := []struct {
		name          string
		code          string
		expectedIssue bool
	}{
		{
			name: "Global variable",
			code: `package test
// This is a global variable
var GlobalConfig = map[string]string{
	"key": "value",
}

func GetConfig() map[string]string {
	return GlobalConfig
}`,
			expectedIssue: true,
		},
		{
			name: "Snake case function",
			code: `package test
func do_something() {
	// Function with snake case name
}`,
			expectedIssue: true,
		},
		{
			name: "Compliant code",
			code: `package test
type Config struct {
	Values map[string]string
}

type Service struct {
	config Config
}

func NewService(config Config) *Service {
	return &Service{config: config}
}

func (s *Service) GetConfig() map[string]string {
	return s.config.Values
}`,
			expectedIssue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzer := NewAnalyzer(AnalyzerConfig{IncludeTests: true})
			file, err := analyzer.ParseString("test.go", tt.code)
			if err != nil {
				t.Fatalf("Failed to parse code: %v", err)
			}

			issues := analyzer.AnalyzeOrganizationStandards(file)
			hasIssue := len(issues) > 0

			if hasIssue != tt.expectedIssue {
				t.Errorf("Expected issue: %v, got: %v", tt.expectedIssue, hasIssue)
				if hasIssue {
					for i, issue := range issues {
						t.Logf("Issue %d: %s at line %d", i+1, issue.Description, issue.Position.Line)
					}
				}
			}
		})
	}
}
