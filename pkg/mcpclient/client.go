package mcpclient

import (
	"encoding/json"
	"fmt"
	"net/rpc/jsonrpc"
)

type Client struct {
	endpoint string
}

type Request struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

type Response struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data,omitempty"`
	Error   string          `json:"error,omitempty"`
}

func New(endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
	}
}

func (c *Client) Call(method string, params interface{}, result interface{}) error {
	conn, err := jsonrpc.Dial("tcp", c.endpoint)
	if err != nil {
		return fmt.Errorf("failed to connect to MCP server: %w", err)
	}
	defer conn.Close()

	// Prepend the service name to the method
	fullMethod := "MCPServer." + method

	var response Response
	if err := conn.Call(fullMethod, params, &response); err != nil {
		return fmt.Errorf("MCP call failed: %w", err)
	}

	if !response.Success {
		return fmt.Errorf("MCP error: %s", response.Error)
	}

	if result != nil && len(response.Data) > 0 {
		if err := json.Unmarshal(response.Data, result); err != nil {
			return fmt.Errorf("failed to unmarshal MCP response: %w", err)
		}
	}

	return nil
}

func (c *Client) GetResource(id string) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := c.Call("GetResource", id, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) GetPrompt(context, fileType, identifier string) (string, error) {
	params := map[string]string{
		"Context":    context,
		"FileType":   fileType,
		"Identifier": identifier,
	}

	var result string
	if err := c.Call("GetPrompt", params, &result); err != nil {
		return "", err
	}
	return result, nil
}

func (c *Client) CallTool(name string, params map[string]interface{}, resourceID string) (interface{}, error) {
	toolRequest := map[string]interface{}{
		"Name":       name,
		"Params":     params,
		"ResourceID": resourceID,
	}

	var result interface{}
	if err := c.Call("CallTool", toolRequest, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) ValidateIntent(content string, ruleIDs []string, fileType string) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"Content":  content,
		"RuleIDs":  ruleIDs,
		"FileType": fileType,
	}

	var result map[string]interface{}
	if err := c.Call("ValidateIntent", params, &result); err != nil {
		return nil, err
	}
	return result, nil
}
