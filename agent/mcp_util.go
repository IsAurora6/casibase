package agent

import (
	"encoding/json"
	"fmt"
)

// 解析JSON并返回第一个服务器的配置参数
func getStdioClientParams(jsonStr string) (string, []string, []string, error) {
	// 解析JSON到map
	var config map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &config); err != nil {
		return "", nil, nil, err
	}

	// 获取mcpServers
	mcpServers, ok := config["mcpServers"].(map[string]interface{})
	if !ok || len(mcpServers) == 0 {
		return "", nil, nil, fmt.Errorf("没有找到服务器配置")
	}

	// 获取第一个服务器配置
	var firstServer map[string]interface{}
	for _, server := range mcpServers {
		firstServer = server.(map[string]interface{})
		break // 只取第一个
	}

	// 提取command
	command, _ := firstServer["command"].(string)

	// 提取args
	argsInterface, _ := firstServer["args"].([]interface{})
	args := make([]string, len(argsInterface))
	for i, arg := range argsInterface {
		args[i] = arg.(string)
	}

	// 提取并转换环境变量
	envMap, _ := firstServer["env"].(map[string]interface{})
	envVars := []string{}
	for k, v := range envMap {
		envVars = append(envVars, fmt.Sprintf("%s=%s", k, v))
	}

	return command, envVars, args, nil
}

// 判断是STDIO还是SSE配置并相应处理
func ParseMCPConfig(jsonStr string) (isSSE bool, command string, envVars []string, args []string, url string, err error) {
	// 先尝试解析为map
	var config map[string]interface{}
	if err = json.Unmarshal([]byte(jsonStr), &config); err != nil {
		return
	}

	// 获取mcpServers
	mcpServers, ok := config["mcpServers"].(map[string]interface{})
	if !ok || len(mcpServers) == 0 {
		err = fmt.Errorf("没有找到服务器配置")
		return
	}

	// 获取第一个服务器配置
	var firstServer map[string]interface{}
	for _, server := range mcpServers {
		firstServer = server.(map[string]interface{})
		break // 只取第一个
	}

	// 检查是否为SSE配置
	if _, hasURL := firstServer["url"]; hasURL {
		isSSE = true
		url, _ = firstServer["url"].(string)
		return
	}

	// 是STDIO配置
	command, envVars, args, err = getStdioClientParams(jsonStr)
	return
}
