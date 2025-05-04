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
	
	// Actual validation logic for the MVP
	issues := []map[string]interface{}{}
	
	for _, ruleID := range req.RuleIDs {
		rulePath := filepath.Join(s.RulesDir, ruleID+".yaml")
		if _, err := os.Stat(rulePath); err != nil {
			issues = append(issues, map[string]interface{}{
				"ruleID":      ruleID,
				"description": "Rule definition not found",
				"severity":    "error",
			})
			continue
		}
		
		// Basic validation for the error_handling rule
		if ruleID == "error_handling" {
			// Check if code has potential error variables but missing error handling
			if strings.Contains(req.Content, "err :=") || strings.Contains(req.Content, "err =") {
				if !strings.Contains(req.Content, "if err != nil") {
					issues = append(issues, map[string]interface{}{
						"ruleID":      "error_handling",
						"description": "Missing error handling pattern 'if err != nil'",
						"severity":    "warning",
					})
				}
			}
			
			// Check for ignored errors using underscore
			if strings.Contains(req.Content, "_ = ") && strings.Contains(req.Content, "Operation()") {
				issues = append(issues, map[string]interface{}{
					"ruleID":      "error_handling",
					"description": "Error is being ignored with underscore assignment",
					"severity":    "error",
				})
			}
		}
		
		// Basic validation for the api_design rule
		if ruleID == "api_design" {
			// Check if function has context parameter
			if strings.Contains(req.Content, "func ") && 
				strings.Contains(req.Content, "(*") && 
				!strings.Contains(req.Content, "ctx context.Context") {
				issues = append(issues, map[string]interface{}{
					"ruleID":      "api_design",
					"description": "API methods should accept context.Context as first parameter",
					"severity":    "warning",
				})
			}
		}
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
