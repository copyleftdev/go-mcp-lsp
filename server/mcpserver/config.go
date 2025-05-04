package mcpserver

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/yourorg/go-mcp-lsp/server/mcpserver/endpoints"
)

type Config struct {
	Address     string `json:"address"`
	RulesDir    string `json:"rulesDir"`
	TemplatesDir string `json:"templatesDir"`
	LogFile     string `json:"logFile,omitempty"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Resolve relative paths
	if !filepath.IsAbs(config.RulesDir) {
		config.RulesDir = filepath.Join(filepath.Dir(path), config.RulesDir)
	}

	if !filepath.IsAbs(config.TemplatesDir) {
		config.TemplatesDir = filepath.Join(filepath.Dir(path), config.TemplatesDir)
	}

	if config.LogFile != "" && !filepath.IsAbs(config.LogFile) {
		config.LogFile = filepath.Join(filepath.Dir(path), config.LogFile)
	}

	return &config, nil
}

func InitializeResourceManager(config *Config) (*endpoints.ResourceManager, error) {
	if _, err := os.Stat(config.RulesDir); err != nil {
		return nil, fmt.Errorf("rules directory not found: %w", err)
	}

	if _, err := os.Stat(config.TemplatesDir); err != nil {
		return nil, fmt.Errorf("templates directory not found: %w", err)
	}

	return endpoints.NewResourceManager(config.RulesDir, config.TemplatesDir), nil
}
