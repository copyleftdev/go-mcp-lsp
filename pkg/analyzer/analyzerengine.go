package analyzer

import (
	"github.com/yourorg/go-mcp-lsp/pkg/analyzer/ast"
)

type AnalysisResult struct {
	Valid  bool
	Issues []ast.Issue
}

type AnalyzerEngine struct {
	analyzer *ast.Analyzer
}

func NewAnalyzerEngine() *AnalyzerEngine {
	config := ast.AnalyzerConfig{
		IncludeTests: false,
	}
	
	return &AnalyzerEngine{
		analyzer: ast.NewAnalyzer(config),
	}
}

func (e *AnalyzerEngine) Analyze(filepath string, content []byte, ruleIDs []string) (*AnalysisResult, error) {
	file, err := e.analyzer.ParseFile(filepath, content)
	if err != nil {
		return nil, err
	}
	
	var allIssues []ast.Issue
	
	for _, ruleID := range ruleIDs {
		var issues []ast.Issue
		
		switch ruleID {
		case "error_handling":
			issues = e.analyzer.AnalyzeErrorHandling(file)
		case "api_design":
			issues = e.analyzer.AnalyzeAPIDesign(file)
		case "concurrent_map_access", "synchronization":
			issues = e.analyzer.AnalyzeConcurrencySafety(file)
		case "secure_coding":
			issues = e.analyzer.AnalyzeSecurityIssues(file)
		case "org_coding_standards", "coding_standards":
			issues = e.analyzer.AnalyzeOrganizationStandards(file)
		}
		
		allIssues = append(allIssues, issues...)
	}
	
	return &AnalysisResult{
		Valid:  len(allIssues) == 0,
		Issues: allIssues,
	}, nil
}
