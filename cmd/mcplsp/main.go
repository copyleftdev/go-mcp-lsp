package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/yourorg/go-mcp-lsp/pkg/analyzer"
	"github.com/yourorg/go-mcp-lsp/pkg/mcpclient"
)

type cliConfig struct {
	mcpEndpoint   string
	rulesDir      string
	mechanismsDir string
	outputFile    string
	command       string
	deep          bool
}

func main() {
	cfg := parseFlags()

	switch cfg.command {
	case "validate":
		validateFile(cfg)
	case "audit":
		auditRules(cfg)
	case "test":
		testConnection(cfg)
	default:
		log.Fatalf("Unknown command: %s", cfg.command)
	}
}

func parseFlags() cliConfig {
	var cfg cliConfig

	flag.StringVar(&cfg.mcpEndpoint, "mcp", "localhost:9000", "MCP server endpoint")
	flag.StringVar(&cfg.rulesDir, "rules", "./pkg/intent", "Path to intent rules")
	flag.StringVar(&cfg.mechanismsDir, "mechanisms", "./pkg/mechanism", "Path to enforcement mechanisms")
	flag.StringVar(&cfg.outputFile, "output", "result.json", "Output file for results")
	flag.BoolVar(&cfg.deep, "deep", false, "Use deep AST-based code inspection (default: false)")

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: mcplsp [flags] <command>")
		fmt.Println("Commands:")
		fmt.Println("  validate <file> [rule1,rule2,...] - Validate a Go file against rules")
		fmt.Println("  audit            - Check for drift between rules and enforcement")
		fmt.Println("  test             - Test connection to MCP server")
		os.Exit(1)
	}

	cfg.command = args[0]
	return cfg
}

func validateFile(cfg cliConfig) {
	if len(flag.Args()) < 2 {
		log.Fatal("Missing file path to validate")
	}

	filePath := flag.Args()[1]
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		log.Fatalf("Failed to resolve file path: %v", err)
	}

	fmt.Printf("Validating file: %s\n", absPath)
	
	content, err := os.ReadFile(absPath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	fmt.Printf("Connecting to MCP server at %s\n", cfg.mcpEndpoint)
	client := mcpclient.New(cfg.mcpEndpoint)
	
	// Get rule IDs to validate against
	ruleIDs := []string{"error_handling", "api_design", "concurrent_map_access", "secure_coding", "org_coding_standards"}
	if len(flag.Args()) > 2 {
		// Use specific rules if provided
		ruleArgs := flag.Args()[2]
		if ruleArgs != "" {
			ruleIDs = strings.Split(ruleArgs, ",")
		}
	}
	
	fmt.Printf("Validating against rules: %v\n", ruleIDs)
	
	if cfg.deep {
		// Use AST-based analyzer for deeper inspection
		fmt.Println("Using deep AST-based code inspection...")
		
		engine := analyzer.NewAnalyzerEngine()
		result, err := engine.Analyze(absPath, content, ruleIDs)
		if err != nil {
			log.Fatalf("Analysis failed: %v", err)
		}
		
		if result.Valid {
			fmt.Println("Validation passed!")
		} else {
			fmt.Println("Validation failed:")
			for _, issue := range result.Issues {
				fmt.Printf("- [%s] %s (Line %d, Col %d) - %s\n", 
					issue.RuleID, 
					issue.Description, 
					issue.Position.Line, 
					issue.Position.Column,
					issue.Severity)
			}
			os.Exit(1)
		}
	} else {
		result, err := client.ValidateCode(string(content), ruleIDs, "go")
		if err != nil {
			log.Fatalf("Validation failed: %v", err)
		}

		fmt.Printf("Raw validation result: %+v\n", result)

		if result.Valid {
			fmt.Println("Validation passed!")
		} else {
			fmt.Println("Validation failed:")
			for _, issue := range result.Issues {
				fmt.Printf("- [%s] %s (Severity: %s)\n", issue.RuleID, issue.Description, issue.Severity)
			}
			os.Exit(1)
		}
	}
}

func auditRules(cfg cliConfig) {
	fmt.Println("Auditing rules for enforcement drift...")
	fmt.Printf("Checking rules in %s against mechanisms in %s\n", cfg.rulesDir, cfg.mechanismsDir)
	fmt.Printf("Results will be written to %s\n", cfg.outputFile)
	
	// In the MVP, we just report that we're auditing without actual implementation
	// This would be expanded in a full implementation
	fmt.Println("MVP: Drift analysis not yet implemented")
}

func testConnection(cfg cliConfig) {
	client := mcpclient.New(cfg.mcpEndpoint)
	
	fmt.Printf("Testing connection to MCP server at %s...\n", cfg.mcpEndpoint)
	
	resource, err := client.GetResource("error_handling")
	if err != nil {
		log.Fatalf("Connection test failed: %v", err)
	}
	
	fmt.Println("Connection successful!")
	fmt.Printf("Retrieved resource: %v\n", resource)
}
