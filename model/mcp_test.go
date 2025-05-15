package model

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

// OAIFunction 封装 OpenAI Function‑Calling 格式
type OAIFunction struct {
	Type     string          `json:"type"` // always "function"
	Function OAIToolFunction `json:"function"`
}
type OAIToolFunction struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Parameters  json.RawMessage `json:"parameters"` // 直接拿 InputSchema 序列化
}

func Test(t *testing.T) {
	ctx := context.Background()
	// 1. 创建 SSE 客户端（需替换为你的高德 Key）
	baseURL := "https://mcp.amap.com/sse?key=69d3b649ca85d0a3175693b6f7444bc3"
	cli, err := client.NewSSEMCPClient(baseURL)
	if err != nil {
		log.Fatalf("创建 SSE MCP 客户端失败: %v", err)
	}
	defer cli.Close()

	// 2. 启动底层连接；必须在 Initialize 之前调用
	if err := cli.Start(ctx); err != nil {
		log.Fatalf("启动 MCP 客户端失败: %v", err)
	}
	// 3. 与服务器协商初始化；必须在调用任何 List* 或 CallTool 之前
	if _, err := cli.Initialize(ctx, mcp.InitializeRequest{}); err != nil {
		log.Fatalf("初始化 MCP 客户端失败: %v", err)
	}

	// 2. 构造并发送 ListTools 请求
	req := mcp.ListToolsRequest{}
	res, err := cli.ListTools(context.Background(), req)
	if err != nil {
		log.Fatalf("调用 ListTools 失败: %v", err)
	}

	//// 3. 打印工具列表
	//fmt.Println("高德 MCP Server 支持以下工具：")
	//for _, tool := range res.Tools {
	//	schemaBytes, _ := json.MarshalIndent(tool.InputSchema, "    ", "  ")
	//	fmt.Printf("- 名称: %s\n", tool.Name)
	//	if tool.Description != "" {
	//		fmt.Printf("  描述: %s\n", tool.Description)
	//	}
	//	fmt.Printf("  输入 Schema Title: %s\n", tool.InputSchema)
	//	fmt.Printf("  完整输入 Schema:\n%s\n\n", string(schemaBytes))
	//}
	// 4. 转换为 OpenAI tools 数组
	var tools []OAIFunction
	for _, t := range res.Tools {
		// 把原始 InputSchema 序列化为 JSON，直接当作 parameters
		paramsJSON, err := json.Marshal(t.InputSchema)
		if err != nil {
			log.Fatalf("序列化 InputSchema 失败: %v", err)
		}
		fn := OAIFunction{
			Type: "function",
			Function: OAIToolFunction{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  paramsJSON,
			},
		}
		tools = append(tools, fn)
	}
	// 5. 输出最终 JSON（可以复制到你的 OpenAI 调用里）
	out, err := json.MarshalIndent(map[string]interface{}{
		"tools": tools,
	}, "", "  ")
	if err != nil {
		log.Fatalf("序列化 tools JSON 失败: %v", err)
	}

	fmt.Println(string(out))
}
