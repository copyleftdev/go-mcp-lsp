package endpoints

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type ResourceManager struct {
	RulesDir     string
	TemplatesDir string
}

type Rule struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Rationale   string   `json:"rationale"`
	Category    string   `json:"category"`
	Severity    string   `json:"severity"`
	Checks      []Check  `json:"checks"`
}

type Check struct {
	Name    string `json:"name"`
	Pattern string `json:"pattern"`
	Ensure  string `json:"ensure"`
	Target  string `json:"target,omitempty"`
	After   string `json:"after,omitempty"`
}

func NewResourceManager(rulesDir, templatesDir string) *ResourceManager {
	return &ResourceManager{
		RulesDir:     rulesDir,
		TemplatesDir: templatesDir,
	}
}

func (rm *ResourceManager) LoadRule(id string) (*Rule, error) {
	path := filepath.Join(rm.RulesDir, id+".yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load rule: %w", err)
	}

	var rule Rule
	if err := json.Unmarshal(data, &rule); err != nil {
		return nil, fmt.Errorf("failed to parse rule: %w", err)
	}

	return &rule, nil
}

func (rm *ResourceManager) ListRules() ([]string, error) {
	var rules []string

	err := filepath.WalkDir(rm.RulesDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(d.Name(), ".yaml") {
			ruleName := strings.TrimSuffix(d.Name(), ".yaml")
			rules = append(rules, ruleName)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list rules: %w", err)
	}

	return rules, nil
}

func (rm *ResourceManager) RenderTemplate(templateName string, data interface{}) (string, error) {
	path := filepath.Join(rm.TemplatesDir, templateName+".tmpl")
	
	tmplData, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to load template: %w", err)
	}
	
	tmpl, err := template.New(templateName).Parse(string(tmplData))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}
	
	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}
	
	return buf.String(), nil
}

func (rm *ResourceManager) ValidateAgainstRule(code string, ruleID string) (bool, []string, error) {
	rule, err := rm.LoadRule(ruleID)
	if err != nil {
		return false, nil, err
	}
	
	var issues []string
	valid := true
	
	// This is a simplified validation that would be replaced with real parsing
	// and AST analysis in a production implementation
	for _, check := range rule.Checks {
		if check.Ensure == "present" {
			if !strings.Contains(code, check.Pattern) {
				issues = append(issues, fmt.Sprintf("Rule %s: Required pattern '%s' not found in code", 
					rule.ID, check.Pattern))
				valid = false
			}
		} else if check.Ensure == "absent" {
			if strings.Contains(code, check.Pattern) {
				issues = append(issues, fmt.Sprintf("Rule %s: Forbidden pattern '%s' found in code", 
					rule.ID, check.Pattern))
				valid = false
			}
		}
	}
	
	return valid, issues, nil
}
