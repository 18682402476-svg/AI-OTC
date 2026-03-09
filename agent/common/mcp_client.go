package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"time"
)

// MCPClient 调用 MCP Server 的客户端
type MCPClient struct {
	BaseURL    string
	AgentToken string
}

// NewMCPClient 初始化 MCP 客户端
func NewMCPClient(token string) *MCPClient {
	return &MCPClient{
		BaseURL:    "http://localhost:8090/api", // 这里应该指向 MCP Server 的网关地址
		AgentToken: token,
	}
}

// 凭证文件结构
type AgentCredentials struct {
	AgentId       int    `json:"agent_id"` // 新增 agent_id 字段
	AgentName     string `json:"agent_name"`
	Token         string `json:"token"`
	WalletAddress string `json:"wallet_address"`
}

const credentialsFile = "agent_credentials.txt"

// saveAgentCredentials 保存凭证到本地文件
func saveAgentCredentials(id int, name, token, walletAddress string) error {
	creds := AgentCredentials{
		AgentId:       id,
		AgentName:     name,
		Token:         token,
		WalletAddress: walletAddress,
	}
	
	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(credentialsFile, data, 0644)
}

// loadAgentCredentials 从本地文件加载凭证
func loadAgentCredentials() (*AgentCredentials, error) {
	// 辅助函数：尝试加载并打印日志
	tryLoad := func(path string) (*AgentCredentials, error) {
		fmt.Printf("尝试加载凭证文件: %s ... ", path)
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("失败: %v\n", err)
			return nil, err
		}
		
		var creds AgentCredentials
		if err := json.Unmarshal(data, &creds); err != nil {
			fmt.Printf("JSON解析失败: %v\n", err)
			// 打印文件内容前50个字符，方便排查
			limit := 50
			if len(data) < limit {
				limit = len(data)
			}
			fmt.Printf("文件内容预览(Hex): %x\n", data[:limit])
			fmt.Printf("文件内容预览(String): %s\n", string(data[:limit]))
			return nil, err
		}
		
		fmt.Printf("成功! Agent: %s (ID: %d)\n", creds.AgentName, creds.AgentId)
		return &creds, nil
	}

	// 1. 尝试当前目录
	if creds, err := tryLoad(credentialsFile); err == nil {
		return creds, nil
	}
	
	// 2. 尝试上级目录 (../agent_credentials.txt)
	if creds, err := tryLoad("../" + credentialsFile); err == nil {
		return creds, nil
	}

	fmt.Println("所有路径均未找到有效的凭证文件，将执行新注册流程。")
	return nil, fmt.Errorf("未找到凭证文件")
}

// LoadCredentials 公开给外部调用，用于获取本地存储的凭证信息
func LoadCredentials() (*AgentCredentials, error) {
	return loadAgentCredentials()
}

// RegisterAgent 在 Java 后端注册 Agent 并获取 Token
// 如果本地已有凭证，则直接复用
func (c *MCPClient) RegisterAgent(agentName, walletAddress string) (string, error) {
	// 1. 尝试从本地加载凭证
	if creds, err := loadAgentCredentials(); err == nil && creds.AgentName == agentName {
		// 简单的校验：如果名字匹配（或者你想复用任何存在的凭证），就直接使用
		// 这里假设如果本地有文件，就优先使用文件里的 Token
		fmt.Printf("发现本地凭证，复用 Agent: %s (ID: %d)\n", creds.AgentName, creds.AgentId)
		c.AgentToken = creds.Token
		return c.AgentToken, nil
	}

	// 2. 本地无凭证，调用注册接口
	payload := map[string]string{
		"agent_name":     agentName,
		"wallet_address": walletAddress,
	}

	respData, err := c.httpPost(fmt.Sprintf("%s/agent/register", c.BaseURL), payload)
	if err != nil {
		return "", err
	}

	// log.Printf("Register response: %s", string(respData)) // Debug

	var result struct {
		Success bool `json:"success"`
		Data    struct {
			AgentId       int    `json:"agentId"` // 解析 agentId
			Token         string `json:"token"`
			WalletAddress string `json:"walletAddress"` 
		} `json:"data"`
		Error string `json:"error"`
	}

	if err := json.Unmarshal([]byte(respData), &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %v, 原始响应: %s", err, string(respData))
	}

	if !result.Success || result.Data.Token == "" {
		msg := result.Error
		if msg == "" {
			msg = "后端未返回 token"
		}
		return "", fmt.Errorf("注册失败: %s", msg)
	}

	// 保存 AgentId, Token 和 WalletAddress 到本地文件
	if err := saveAgentCredentials(result.Data.AgentId, agentName, result.Data.Token, result.Data.WalletAddress); err != nil {
		fmt.Printf("警告: 无法保存 Agent 凭证: %v\n", err)
	}

	c.AgentToken = result.Data.Token
	return c.AgentToken, nil
}

// FetchAllMarketData 聚合来自 Java 后端的所有市场信息
func (c *MCPClient) FetchAllMarketData(tradingPair, tokenAddress string) (*MarketDataContext, error) {
	rate, err := c.httpGet(fmt.Sprintf("%s/rate/current?tradingPair=%s", c.BaseURL, tradingPair))
	if err != nil {
		return nil, err
	}

	// 动态计算最近 24 小时的时间范围
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	layout := "2006-01-02 15:04:05"
	historyUrl := fmt.Sprintf("%s/rate/history?tradingPair=%s&startTime=%s&endTime=%s", 
		c.BaseURL, tradingPair, yesterday.Format(layout), now.Format(layout))

	history, err := c.httpGet(historyUrl)
	if err != nil {
		return nil, err
	}

	orders, err := c.httpGetAuth(fmt.Sprintf("%s/order/active", c.BaseURL))
	if err != nil {
		return nil, err
	}

	myOrders, err := c.httpGetAuth(fmt.Sprintf("%s/order/list", c.BaseURL))
	if err != nil {
		return nil, err
	}

	// 1. 获取本地凭证中的 Agent ID
	creds, err := loadAgentCredentials()
	if err != nil {
		return nil, fmt.Errorf("无法加载本地凭证获取 Agent ID: %v", err)
	}

	// 2. 调用 /agent/detail/{agentId} 接口获取余额
	detailUrl := fmt.Sprintf("%s/agent/detail/%d", c.BaseURL, creds.AgentId)
	agentDetail, err := c.httpGetAuth(detailUrl)
	if err != nil {
		return nil, err
	}
	
	// 打印余额详情，方便调试
	fmt.Printf("Agent Detail Response: %s\n", agentDetail)

	return &MarketDataContext{
		CurrentPrice: rate,
		History24h:   history,
		ActiveOrders: orders,
		MyOrders:     myOrders,
		MyBalance:    agentDetail, // 这里使用 Agent Detail 作为余额信息
	}, nil
}

// ExecuteAction 根据 LLM 的决策执行真实的交易操作
// 增加 currentRate 参数，用于计算总价
func (c *MCPClient) ExecuteAction(action *AgentAction, currentRate float64) (string, error) {
	switch action.Action {
	case "HOLD":
		return "保持不动", nil
	case "SELL_MON_FOR_USDC":
		// 卖出 MON: mon_amount 是数量
		// usdc_expected 是 LLM 预期的总价
		amount := toFloat64(action.Params["mon_amount"])
		
		// 优先使用 LLM 预期的总价，如果未提供或为0，则使用实时汇率计算
		totalUSDC := toFloat64(action.Params["usdc_expected"])
		if currentRate > 0 {
			totalUSDC = amount * currentRate
			fmt.Printf("LLM 未提供 usdc_expected，使用实时汇率 %f 计算总价: %f\n", currentRate, totalUSDC)
		}
		
		// 用户坚持 Price 字段是总价
		price := totalUSDC

		// 精度处理
		amount = math.Floor(amount*1000000) / 1000000
		price = math.Floor(price*1000000) / 1000000
		
		// 再次强制确保是浮点数
		if amount == math.Floor(amount) {
			amount += 0.000001
		}
		if price == math.Floor(price) {
			price += 0.000001
		}
		
		// 获取 MON 的地址
		monAddress := os.Getenv("MON_TOKEN_ADDRESS")
		if monAddress == "" {
			monAddress = "0xf8829110cab77c895b1545965DcF34d797dBA295"
		}
		
		return c.httpPostAuth(fmt.Sprintf("%s/order/create", c.BaseURL), map[string]interface{}{
			"type":           "SELL", // SELL_MON_FOR_USDC -> SELL (卖 MON)
			"tokenAddress":   monAddress, // 补全 MON 地址
			"tokenSymbol":    "MON",
			"amount":         amount,
			"price":          price, // 发送总价
			"thoughtProcess": action.Reasoning,
		})
	case "SELL_USDC_FOR_MON":
		// 买入 MON: usdc_amount 是总花费(USDC)，mon_expected 是期望买到的数量(MON)
		totalUSDC := toFloat64(action.Params["usdc_amount"])
		expectedMON := toFloat64(action.Params["mon_expected"])
		
		// 这里的 amount 指的是要买入的 MON 数量
		amount := expectedMON
		
		// 如果 LLM 没给总花费，用 amount * 汇率 估算
		if amount > 0 && currentRate > 0 {
			totalUSDC = amount * currentRate
			fmt.Printf("LLM 未提供 usdc_amount，使用实时汇率 %f 计算总价: %f\n", currentRate, totalUSDC)
		}
		
		// 用户坚持 Price 字段是总价
		price := totalUSDC

		// 精度处理
		amount = math.Floor(amount*1000000) / 1000000
		price = math.Floor(price*1000000) / 1000000

		// 再次强制确保是浮点数
		if amount == math.Floor(amount) {
			amount += 0.000001
		}
		if price == math.Floor(price) {
			price += 0.000001
		}

		// 即使是 BUY (买 MON)，后端也要求 tokenAddress 必须是 MON 的地址
		monAddress := os.Getenv("MON_TOKEN_ADDRESS")
		if monAddress == "" {
			monAddress = "0xf8829110cab77c895b1545965DcF34d797dBA295"
		}

		return c.httpPostAuth(fmt.Sprintf("%s/order/create", c.BaseURL), map[string]interface{}{
			"type":           "BUY", // SELL_USDC_FOR_MON -> BUY (买 MON)
			"tokenAddress":   monAddress, // 这里必须填 MON 地址，而不是 USDC 地址
			"tokenSymbol":    "MON",      // 这里的 Symbol 也应该是 MON
			"amount":         amount,
			"price":          price, // 发送总价
			"thoughtProcess": action.Reasoning,
		})
	case "BUY_MON_FROM_ORDER":
		return c.httpPostAuth(fmt.Sprintf("%s/order/buy", c.BaseURL), map[string]interface{}{
			"orderId":        action.Params["order_id"],
			"thoughtProcess": action.Reasoning,
		})
	case "BUY_USDC_FROM_ORDER":
		return c.httpPostAuth(fmt.Sprintf("%s/order/sell", c.BaseURL), map[string]interface{}{
			"orderId":        action.Params["order_id"],
			"thoughtProcess": action.Reasoning,
		})
	case "RE_ORDER":
		return c.httpPostAuth(fmt.Sprintf("%s/order/cancel", c.BaseURL), map[string]interface{}{
			"orderId":        action.Params["order_id"],
			"thoughtProcess": action.Reasoning,
		})
	default:
		return "", fmt.Errorf("未知动作类型: %s", action.Action)
	}
}

func (c *MCPClient) httpGet(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (c *MCPClient) httpGetAuth(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	if c.AgentToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AgentToken)
	}
	if k := os.Getenv("MCP_API_KEY"); k != "" {
		req.Header.Set("X-API-Key", k)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (c *MCPClient) httpPost(url string, payload interface{}) (string, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// toFloat64 尝试将 interface{} 转换为 float64
// 解决 Java 后端 Integer cannot cast to Double 的问题
func toFloat64(v interface{}) float64 {
	if v == nil {
		return 0.0
	}
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case string:
		// 尝试从字符串解析
		var f float64
		fmt.Sscanf(val, "%f", &f)
		return f
	default:
		return 0.0
	}
}

// ParseJSON 辅助函数，用于解析 JSON 字符串
func ParseJSON(jsonStr string, v interface{}) error {
	return json.Unmarshal([]byte(jsonStr), v)
}

func (c *MCPClient) httpPostAuth(url string, payload interface{}) (string, error) {
	var jsonData []byte
	var err error
	if payload != nil {
		jsonData, err = json.Marshal(payload)
		if err != nil {
			return "", err
		}
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.AgentToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AgentToken)
	}
	if k := os.Getenv("MCP_API_KEY"); k != "" {
		req.Header.Set("X-API-Key", k)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
