package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"strings"
)

type Config struct {
	JavaBackendURL string `json:"java_backend_url"`
}

func loadConfig() string {
	// 默认地址
	// defaultURL := "http://192.168.12.51:8080"
	defaultURL := "http://127.0.0.1:8080"

	// 读取配置文件
	data, err := os.ReadFile("config.json")
	if err != nil {
		log.Printf("未找到配置文件 config.json，使用默认地址: %s", defaultURL)
		return defaultURL
	}

	// 检查是否被 DLP 软件加密 (以 %TSD 开头)
	if len(data) > 4 && strings.HasPrefix(string(data), "%TSD") {
		log.Printf("检测到配置文件被加密 (%s)，忽略文件内容，使用默认地址: %s", string(data[:15]), defaultURL)
		return defaultURL
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		log.Printf("配置文件解析失败: %v，使用默认地址: %s", err, defaultURL)
		// 打印文件内容预览 (Hex)，用于调试文件损坏问题
		limit := 50
		if len(data) < limit {
			limit = len(data)
		}
		log.Printf("文件内容预览(Hex): %x", data[:limit])
		log.Printf("文件内容预览(String): %s", string(data[:limit]))
		return defaultURL
	}

	if config.JavaBackendURL == "" {
		log.Printf("配置文件中未指定 java_backend_url，使用默认地址: %s", defaultURL)
		return defaultURL
	}

	log.Printf("从配置文件加载后端地址: %s", config.JavaBackendURL)
	return config.JavaBackendURL
}

func main() {
	// 创建一个新的 MCP 服务器实例
	s := server.NewMCPServer(
		"AgentOtcPlatform-MCP", // 服务名称
		"1.0.0",                 // 服务版本
		server.WithLogging(),    // 启用日志记录
	)

	// 从配置文件加载 Java 后端地址
	javaBaseURL := loadConfig()

	// 启动 HTTP 网关，监听 8090，将 /api/** 转发到 Java 后端
	startHTTPGateway(javaBaseURL)

	// 注册各类工具 (Tools)
	registerMarketTools(s, javaBaseURL) // 市场行情相关工具
	registerOrderTools(s, javaBaseURL)  // 订单管理相关工具
	registerAgentTools(s, javaBaseURL)  // Agent 个人资产相关工具

	// 启动服务器，使用标准输入输出 (stdio) 进行通信
	// 这是 MCP Server 与 AI 客户端（如 Claude Desktop）通信的标准方式
	// if err := server.ServeStdio(s); err != nil {
	// 	log.Fatalf("MCP 服务器启动失败: %v", err)
	// }

	// 阻塞主进程，保持 HTTP 网关运行
	select {}
}

// registerMarketTools 注册市场行情相关的工具
func registerMarketTools(s *server.MCPServer, baseURL string) {
	// 定义 get_market_rates 工具：获取实时汇率及 24 小时历史趋势
	tool := mcp.NewTool("get_market_rates",
		mcp.WithDescription("获取交易对的当前市场汇率及 24 小时历史记录"),
	)

	// 实现工具的具体逻辑
	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		tradingPair, ok := request.Params.Arguments["tradingPair"].(string)
		if !ok {
			return mcp.NewToolResultError("tradingPair 参数是必填的"), nil
		}

		// 调用 api_handlers.go 中的逻辑访问 Java 后端
		return handleGetMarketRates(baseURL, tradingPair)
	})
}

// registerOrderTools 注册订单相关的工具
func registerOrderTools(s *server.MCPServer, baseURL string) {
	// 定义 list_platform_orders 工具：列出平台上所有活跃的挂单
	tool := mcp.NewTool("list_platform_orders",
		mcp.WithDescription("查询平台上当前所有的活跃挂单列表"),
	)

	s.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return handleListPlatformOrders(baseURL)
	})
}

// registerAgentTools 注册 Agent 个人相关的工具
func registerAgentTools(s *server.MCPServer, baseURL string) {
	// 定义 get_agent_orders 工具：查询特定 Agent 的订单记录
	toolOrders := mcp.NewTool("get_agent_orders",
		mcp.WithDescription("查询特定 Agent 的历史订单及当前挂单"),
	)

	s.AddTool(toolOrders, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		token, _ := request.Params.Arguments["token"].(string)
		return handleGetAgentOrders(baseURL, token)
	})

	// 定义 get_agent_balance 工具：查询 Agent 的钱包余额
	toolBalance := mcp.NewTool("get_agent_balance",
		mcp.WithDescription("查询特定 Agent 在平台上的代币余额"),
	)

	s.AddTool(toolBalance, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		token, _ := request.Params.Arguments["token"].(string)
		tokenAddress, _ := request.Params.Arguments["tokenAddress"].(string)
		return handleGetAgentBalance(baseURL, token, tokenAddress)
	})
}
