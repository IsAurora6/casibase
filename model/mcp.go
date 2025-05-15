package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// MCPServerConfig 是 MCP 服务器配置
type MCPServerConfig struct {
	URL string `json:"url"`
}

// MCPServersConfig 是所有 MCP 服务器的配置
type MCPServersConfig struct {
	MCPServers map[string]MCPServerConfig `json:"mcpServers"`
}

// MCPTool 表示 MCP 服务提供的工具
type MCPTool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
}

// 从 MCP 服务器获取工具列表
func getMCPTools(mcpServerURL string) ([]MCPTool, error) {
	req, err := http.NewRequest("GET", mcpServerURL+"/tools", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call MCP server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("MCP server returned error status: %d", resp.StatusCode)
	}

	var toolsResponse struct {
		Tools []MCPTool `json:"tools"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&toolsResponse); err != nil {
		return nil, fmt.Errorf("failed to decode tools response: %v", err)
	}

	return toolsResponse.Tools, nil
}

// 调用 MCP 工具
func callMCPTool(mcpServerURL, toolName string, parameters map[string]interface{}) (interface{}, error) {
	toolCallReq := struct {
		Name       string                 `json:"name"`
		Parameters map[string]interface{} `json:"parameters"`
	}{
		Name:       toolName,
		Parameters: parameters,
	}

	reqBody, err := json.Marshal(toolCallReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tool call request: %v", err)
	}

	req, err := http.NewRequest("POST", mcpServerURL+"/call_tool", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call MCP server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("MCP server returned error status: %d", resp.StatusCode)
	}

	var toolResponse struct {
		Status string      `json:"status"`
		Result interface{} `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&toolResponse); err != nil {
		return nil, fmt.Errorf("failed to decode tool call response: %v", err)
	}

	return toolResponse.Result, nil
}
