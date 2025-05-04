package mcpclient

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Endpoint string `json:"endpoint"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if config.Endpoint == "" {
		return nil, fmt.Errorf("missing endpoint in config")
	}

	return &config, nil
}

func NewFromConfig(configPath string) (*Client, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}
	
	return New(config.Endpoint), nil
}

type ValidationResult struct {
	Valid   bool     `json:"valid"`
	Issues  []Issue  `json:"issues,omitempty"`
}

type Issue struct {
	RuleID      string `json:"ruleID"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
	Location    *Location `json:"location,omitempty"`
}

type Location struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

func (c *Client) ValidateCode(content string, ruleIDs []string, fileType string) (*ValidationResult, error) {
	result, err := c.ValidateIntent(content, ruleIDs, fileType)
	if err != nil {
		return nil, err
	}
	
	validationBytes, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal validation result: %w", err)
	}
	
	var validationResult ValidationResult
	if err := json.Unmarshal(validationBytes, &validationResult); err != nil {
		return nil, fmt.Errorf("failed to unmarshal validation result: %w", err)
	}
	
	return &validationResult, nil
}
