package mcpserver

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"path/filepath"
	"strings"
	
	"github.com/yourorg/go-mcp-lsp/pkg/analyzer"
)

type MCPServer struct {
	RulesDir     string
	TemplatesDir string
	listener     net.Listener
}

type Resource struct {
	ID      string
	Content string
	Type    string
}

type PromptRequest struct {
	Context    string
	FileType   string
	Identifier string
}

type ToolRequest struct {
	Name       string
	Params     map[string]interface{}
	ResourceID string
}

type ValidateRequest struct {
	Content  string
	RuleIDs  []string
	FileType string
}

type Result struct {
	Success bool
	Data    interface{}
	Error   string
}

func NewMCPServer(rulesDir, templatesDir string) (*MCPServer, error) {
	if _, err := os.Stat(rulesDir); err != nil {
		return nil, fmt.Errorf("rules directory not found: %w", err)
	}
	
	if _, err := os.Stat(templatesDir); err != nil {
		return nil, fmt.Errorf("templates directory not found: %w", err)
	}
	
	return &MCPServer{
		RulesDir:     rulesDir,
		TemplatesDir: templatesDir,
	}, nil
}

func (s *MCPServer) Start(address string) error {
	rpc.Register(s)
	
	var err error
	s.listener, err = net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to start MCP server: %w", err)
	}
	
	log.Printf("MCP server listening on %s", address)
	
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
			continue
		}
		
		go jsonrpc.ServeConn(conn)
	}
}

func (s *MCPServer) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *MCPServer) GetResource(id string, result *Result) error {
	filePath := filepath.Join(s.RulesDir, id+".yaml")
	
	data, err := os.ReadFile(filePath)
	if err != nil {
		*result = Result{
			Success: false,
			Error:   fmt.Sprintf("resource not found: %v", err),
		}
		return nil
	}
	
	*result = Result{
		Success: true,
		Data: Resource{
			ID:      id,
			Content: string(data),
			Type:    "yaml",
		},
	}
	
	return nil
}

func (s *MCPServer) GetPrompt(req PromptRequest, result *Result) error {
	templatePath := filepath.Join(s.TemplatesDir, req.FileType, req.Identifier+".tmpl")
	
	data, err := os.ReadFile(templatePath)
	if err != nil {
		*result = Result{
			Success: false,
			Error:   fmt.Sprintf("template not found: %v", err),
		}
		return nil
	}
	
	*result = Result{
		Success: true,
		Data:    string(data),
	}
	
	return nil
}

func (s *MCPServer) CallTool(req ToolRequest, result *Result) error {
	switch req.Name {
	case "generateScaffold":
		if template, ok := req.Params["template"].(string); ok {
			templatePath := filepath.Join(s.TemplatesDir, template+".tmpl")
			data, err := os.ReadFile(templatePath)
			if err != nil {
				*result = Result{
					Success: false,
					Error:   fmt.Sprintf("template not found: %v", err),
				}
				return nil
			}
			
			*result = Result{
				Success: true,
				Data:    string(data),
			}
			return nil
		}
		*result = Result{
			Success: false,
			Error:   "missing template parameter",
		}
		return nil

	case "validateCode":
		if code, ok := req.Params["code"].(string); ok {
			// Simple validation logic - in a real implementation, analyze the code
			valid := true
			var messages []string
			
			// Simple check for error handling (example only)
			if strings.Contains(code, "if err != nil") {
				valid = true
			} else if strings.Contains(code, "err :=") || strings.Contains(code, "err =") {
				valid = false
				messages = append(messages, "Missing error handling")
			}
			
			*result = Result{
				Success: true,
				Data: map[string]interface{}{
					"valid": valid,
					"messages": messages,
				},
			}
			return nil
		}
		*result = Result{
			Success: false,
			Error:   "missing code parameter",
		}
		return nil

	default:
		*result = Result{
			Success: false,
			Error:   fmt.Sprintf("unknown tool: %s", req.Name),
		}
		return nil
	}
}

func (s *MCPServer) ValidateIntent(req ValidateRequest, result *Result) error {
	if len(req.RuleIDs) == 0 {
		*result = Result{
			Success: false,
			Error:   "no rule IDs specified",
		}
		return nil
	}
	
	// Use AST-based analyzer for deeper code inspection
	engine := analyzer.NewAnalyzerEngine()
	analysisResult, err := engine.Analyze("file.go", []byte(req.Content), req.RuleIDs)
	
	if err != nil {
		*result = Result{
			Success: false,
			Error:   fmt.Sprintf("analysis failed: %v", err),
		}
		return nil
	}
	
	// Map analyzer issues to MCP issues format
	issues := []map[string]interface{}{}
	
	if analysisResult != nil && len(analysisResult.Issues) > 0 {
		for _, issue := range analysisResult.Issues {
			issueMap := map[string]interface{}{
				"ruleID":      issue.RuleID,
				"description": issue.Description,
				"severity":    issue.Severity,
				"location": map[string]interface{}{
					"line":   issue.Position.Line,
					"column": issue.Position.Column,
				},
			}
			issues = append(issues, issueMap)
		}
	} else {
		// Fallback to simpler pattern-based checks if AST analysis doesn't find issues
		// This can help catch issues that might be missed by AST analysis
		issues = performPatternBasedValidation(req.Content, req.RuleIDs)
	}
	
	*result = Result{
		Success: true,
		Data: map[string]interface{}{
			"valid":  len(issues) == 0,
			"issues": issues,
		},
	}
	
	return nil
}

// performPatternBasedValidation implements the previous string-based validation logic
// as a fallback mechanism when AST analysis doesn't find any issues
func performPatternBasedValidation(content string, ruleIDs []string) []map[string]interface{} {
	issues := []map[string]interface{}{}
	
	for _, ruleID := range ruleIDs {
		// Handle subdirectory paths in rule IDs
		ruleActualID := ruleID
		
		// If the rule ID contains a path separator, extract the actual rule ID
		if strings.Contains(ruleID, "/") {
			parts := strings.Split(ruleID, "/")
			ruleActualID = parts[len(parts)-1]
		}
		
		// Basic validation based on the rule ID
		switch ruleActualID {
		case "error_handling":
			// Check for error handling patterns
			if strings.Contains(content, "err :=") || strings.Contains(content, "err =") {
				if !strings.Contains(content, "if err != nil") {
					issues = append(issues, map[string]interface{}{
						"ruleID":      "error_handling",
						"description": "Missing error handling pattern 'if err != nil'",
						"severity":    "warning",
					})
				}
			}
			
			// Check for ignored errors
			if strings.Contains(content, "_ =") && strings.Contains(content, "err") {
				issues = append(issues, map[string]interface{}{
					"ruleID":      "error_handling",
					"description": "Error is being ignored with underscore assignment",
					"severity":    "error",
				})
			}
			
		case "api_design":
			// Check for context parameter
			if strings.Contains(content, "func ") && 
				strings.Contains(content, "(*") && 
				!strings.Contains(content, "ctx context.Context") {
				issues = append(issues, map[string]interface{}{
					"ruleID":      "api_design",
					"description": "API methods should accept context.Context as first parameter",
					"severity":    "warning",
				})
			}
			
		case "concurrent_map_access", "synchronization":
			// Check for concurrent map access without synchronization
			if strings.Contains(content, "go func") &&
				strings.Contains(content, "map[") &&
				!strings.Contains(content, "sync.Mutex") &&
				!strings.Contains(content, "sync.RWMutex") {
				issues = append(issues, map[string]interface{}{
					"ruleID":      "concurrent_map_access",
					"description": "Concurrent map access without proper synchronization",
					"severity":    "error",
				})
			}
			
		case "secure_coding":
			// Check for weak crypto
			if strings.Contains(content, "crypto/md5") || strings.Contains(content, "crypto/sha1") {
				issues = append(issues, map[string]interface{}{
					"ruleID":      "secure_coding",
					"description": "Using weak cryptographic algorithms (MD5/SHA1)",
					"severity":    "error",
				})
			}
			
			// Check for potential SQL injection
			if strings.Contains(content, "fmt.Sprintf") && 
				strings.Contains(content, "SELECT") &&
				strings.Contains(content, "%s") {
				issues = append(issues, map[string]interface{}{
					"ruleID":      "secure_coding",
					"description": "Potential SQL injection vulnerability",
					"severity":    "error",
				})
			}
			
			// Check for hardcoded credentials
			credPatterns := []string{"password :=", "apiKey :=", "secret :=", "token :="}
			for _, pattern := range credPatterns {
				if strings.Contains(content, pattern) && 
					strings.Contains(content, "\"") &&
					!strings.Contains(content, "os.Getenv") {
					issues = append(issues, map[string]interface{}{
						"ruleID":      "secure_coding",
						"description": "Hardcoded credentials detected",
						"severity":    "error",
					})
					break
				}
			}
			
		case "coding_standards", "org_coding_standards":
			// Check for global variables
			if strings.Contains(content, "var ") &&
				strings.Contains(content, "Global") {
				issues = append(issues, map[string]interface{}{
					"ruleID":      "org_coding_standards",
					"description": "Global variables violate organizational standards",
					"severity":    "warning",
				})
			}
			
			// Check for snake_case function names
			if strings.Contains(content, "func do_") || strings.Contains(content, "func get_") {
				issues = append(issues, map[string]interface{}{
					"ruleID":      "org_coding_standards",
					"description": "Snake case function names are not allowed",
					"severity":    "warning",
				})
			}
			
			// Check for dependency injection patterns
			if strings.Contains(content, "type Service struct") &&
				!strings.Contains(content, "Config struct") {
				issues = append(issues, map[string]interface{}{
					"ruleID":      "org_coding_standards",
					"description": "Missing configuration struct for dependency injection",
					"severity":    "warning",
				})
			}
		}
	}
	
	return issues
}
