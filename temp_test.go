package main

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/casibase/casibase/agent"

	"github.com/mark3labs/mcp-go/client"
	mcpgo "github.com/mark3labs/mcp-go/mcp"
)

func Test(t *testing.T) {

	ctx := context.Background()

	configStdio := `{
			"mcpServers":{
				"baidu-map": {
					"command": "npx",
					"args": [
						"-y",
						"@baidumap/agent-server-baidu-map"
					],
					"env": {
						"BAIDU_MAP_API_KEY": "oIs9jLfPLrbCDUE8Pc6Bkz1S6Vl6SF2W"
					}
				}
			}
		  }`
	//configSSE := `{
	//		"mcpServers": {
	//			"amap-amap-sse": {
	//  				"url": "https://mcp.amap.com/sse?key=69d3b649ca85d0a3175693b6f7444bc3"
	//			}
	//	  	}
	//	}`
	// Create client based on transport type
	var c *client.Client
	var err error
	isSSE, command, envVars, args, url, err := agent.ParseMCPConfig(configStdio)
	if !isSSE {
		c, err = client.NewStdioMCPClient(
			command,
			envVars,
			args...,
		)
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}
	} else {
		// 1. 创建 SSE 客户端（需替换为你的高德 Key）
		c, err = client.NewSSEMCPClient(url)
		if err != nil {
			log.Fatalf("创建 SSE MCP 客户端失败: %v", err)
		}
	}

	// Initialize the client
	if err := c.Start(ctx); err != nil {
		log.Fatalf("启动 MCP 客户端失败: %v", err)
	}
	// 3. 与服务器协商初始化；必须在调用任何 List* 或 CallTool 之前
	serverInfo, err := c.Initialize(ctx, mcpgo.InitializeRequest{})
	if err != nil {
		log.Fatalf("初始化 MCP 客户端失败: %v", err)
	}

	// Display server information
	fmt.Printf("Connected to server: %s (version %s)\n",
		serverInfo.ServerInfo.Name,
		serverInfo.ServerInfo.Version)
	fmt.Printf("Server capabilities: %+v\n", serverInfo.Capabilities)

	// List available tools if the server supports them
	if serverInfo.Capabilities.Tools != nil {
		fmt.Println("Fetching available tools...")
		toolsRequest := mcpgo.ListToolsRequest{}
		toolsResult, err := c.ListTools(ctx, toolsRequest)
		if err != nil {
			log.Printf("Failed to list tools: %v", err)
		} else {
			fmt.Printf("Server has %d tools available\n", len(toolsResult.Tools))
			for i, tool := range toolsResult.Tools {
				fmt.Printf("  %d. %s - %s\n", i+1, tool.Name, tool.Description)
			}
		}
	}

	// List available resources if the server supports them
	if serverInfo.Capabilities.Resources != nil {
		fmt.Println("Fetching available resources...")
		resourcesRequest := mcpgo.ListResourcesRequest{}
		resourcesResult, err := c.ListResources(ctx, resourcesRequest)
		if err != nil {
			log.Printf("Failed to list resources: %v", err)
		} else {
			fmt.Printf("Server has %d resources available\n", len(resourcesResult.Resources))
			for i, resource := range resourcesResult.Resources {
				fmt.Printf("  %d. %s - %s\n", i+1, resource.URI, resource.Name)
			}
		}
	}

	fmt.Println("Client initialized successfully. Shutting down...")
	c.Close()
}
