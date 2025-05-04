package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/yourorg/go-mcp-lsp/pkg/analyzer/ast"
)

func main() {
	analyzeCmd := flag.NewFlagSet("analyze", flag.ExitOnError)
	filePath := analyzeCmd.String("file", "", "Path to the Go file to analyze")
	rulesStr := analyzeCmd.String("rules", "", "Comma-separated list of rules to check")

	if len(os.Args) < 2 {
		fmt.Println("Expected 'analyze' subcommand")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "analyze":
		analyzeCmd.Parse(os.Args[2:])
		if *filePath == "" {
			fmt.Println("--file flag is required")
			os.Exit(1)
		}
		if *rulesStr == "" {
			fmt.Println("--rules flag is required")
			os.Exit(1)
		}
		runAnalysis(*filePath, *rulesStr)
	default:
		fmt.Printf("Unknown subcommand: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func runAnalysis(filePath, rulesStr string) {
	// Read the file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Initialize the analyzer
	config := ast.AnalyzerConfig{
		IncludeTests: true,
	}
	analyzer := ast.NewAnalyzer(config)

	// Parse the file
	file, err := analyzer.ParseFile(filePath, content)
	if err != nil {
		fmt.Printf("Error parsing file: %v\n", err)
		os.Exit(1)
	}

	// Split the rules string
	rules := strings.Split(rulesStr, ",")

	// Run the appropriate analysis for each rule
	for _, rule := range rules {
		rule = strings.TrimSpace(rule)
		var issues []ast.Issue

		switch rule {
		case "error_handling":
			issues = analyzer.AnalyzeErrorHandling(file)
		case "api_design":
			issues = analyzer.AnalyzeAPIDesign(file)
		case "concurrent_map_access", "synchronization":
			issues = analyzer.AnalyzeConcurrencySafety(file)
		case "secure_coding":
			issues = analyzer.AnalyzeSecurityIssues(file)
		case "org_coding_standards", "coding_standards":
			issues = analyzer.AnalyzeOrganizationStandards(file)
		default:
			fmt.Printf("Unknown rule: %s\n", rule)
			continue
		}

		// Print the results
		fmt.Printf("Rule: %s\n", rule)
		if len(issues) == 0 {
			fmt.Println("  No issues found")
		} else {
			for i, issue := range issues {
				fmt.Printf("  Issue %d:\n", i+1)
				fmt.Printf("    Description: %s\n", issue.Description)
				fmt.Printf("    Severity: %s\n", issue.Severity)
				fmt.Printf("    Location: Line %d, Column %d\n", issue.Position.Line, issue.Position.Column)
			}
		}
	}
}
