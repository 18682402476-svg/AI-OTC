package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

// handleGetMarketRates 调用 Java 后端的 /rate/current 和 /rate/history 接口
// 聚合实时价格和 24 小时历史趋势数据
func handleGetMarketRates(baseURL, tradingPair string) (*mcp.CallToolResult, error) {
	// 1. 获取当前最新汇率
	currentURL := fmt.Sprintf("%s/api/rate/current?tradingPair=%s", baseURL, tradingPair)
	// log.Printf("Fetching market rates from: %s", currentURL) // Debug log
	currentData, err := httpGet(currentURL, "")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("获取实时汇率失败: %v", err)), nil
	}

	// 2. 获取历史记录（最近 24 小时）
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	layout := "2006-01-02 15:04:05"
	historyURL := fmt.Sprintf("%s/api/rate/history?tradingPair=%s&startTime=%s&endTime=%s",
		baseURL, tradingPair, yesterday.Format(layout), now.Format(layout))
	// log.Printf("Fetching history rates from: %s", historyURL) // Debug log
	historyData, err := httpGet(historyURL, "")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("获取汇率历史失败: %v", err)), nil
	}

	// 聚合两组数据返回给 Agent
	result := map[string]interface{}{
		"current": currentData["data"],
		"history": historyData["data"],
	}

	jsonBytes, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// handleListPlatformOrders 调用 /order/active 接口
// 返回平台上所有正在进行的挂单，供 Agent 寻找交易对手
func handleListPlatformOrders(baseURL string) (*mcp.CallToolResult, error) {
	url := fmt.Sprintf("%s/api/order/active", baseURL)
	// log.Printf("Fetching active orders from: %s", url) // Debug log
	// 注意：如果 Java 后端启用了全局拦截器，这里可能需要传入一个默认的管理员 Token 或让该接口支持匿名访问
	data, err := httpGet(url, "")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("获取平台挂单失败: %v", err)), nil
	}

	jsonBytes, _ := json.MarshalIndent(data["data"], "", "  ")
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// handleGetAgentOrders 调用 /order/list 接口
// 根据 Agent 的 Token 查询其个人的所有订单状态
func handleGetAgentOrders(baseURL, token string) (*mcp.CallToolResult, error) {
	url := fmt.Sprintf("%s/api/order/list", baseURL)
	// log.Printf("Fetching agent orders from: %s", url) // Debug log
	data, err := httpGet(url, token)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("获取 Agent 订单失败: %v", err)), nil
	}

	jsonBytes, _ := json.MarshalIndent(data["data"], "", "  ")
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// handleGetAgentBalance 调用 /fund/balance 接口
// 查询 Agent 子账户中特定代币的可用余额
func handleGetAgentBalance(baseURL, token, tokenAddress string) (*mcp.CallToolResult, error) {
	url := fmt.Sprintf("%s/api/fund/balance?tokenAddress=%s", baseURL, tokenAddress)
	// log.Printf("Fetching agent balance from: %s", url) // Debug log
	data, err := httpGet(url, token)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("获取 Agent 余额失败: %v", err)), nil
	}

	jsonBytes, _ := json.MarshalIndent(data["data"], "", "  ")
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// httpGet 通用的 HTTP GET 请求辅助函数
// 支持在 Header 中携带 Bearer Token 进行身份鉴权
func httpGet(url string, token string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 如果提供了 Token，则添加到 Authorization 头部
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	// 设置 10 秒超时，防止请求挂死
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应体内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 检查 HTTP 状态码是否为 200
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("服务器返回错误状态码 %d: %s", resp.StatusCode, string(body))
	}

	// 将 JSON 响应解析为 Go map
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func httpPost(url string, token string, payload []byte) ([]byte, int, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		// 检查 token 是否已经包含 "Bearer " 前缀，避免重复添加
		if !strings.HasPrefix(token, "Bearer ") {
			req.Header.Set("Authorization", "Bearer "+token)
		} else {
			req.Header.Set("Authorization", token)
		}
	}
	// 增加超时时间到 120 秒，因为区块链交易可能需要较长时间确认
	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return body, resp.StatusCode, nil
}

func startHTTPGateway(javaBaseURL string) {
	mux := http.NewServeMux()
	rl := newRateLimiter(120, time.Minute)
	mux.HandleFunc("/api/rate/current", func(w http.ResponseWriter, r *http.Request) {
		if !checkAPIKey(w, r) || !rl.allow(r.RemoteAddr) {
			return
		}
		log.Printf("GET %s %s", r.URL.Path, r.URL.RawQuery)
		target := javaBaseURL + "/api/rate/current?" + r.URL.RawQuery
		data, code, err := simpleProxyGet(target, r.Header.Get("Authorization"))
		if err != nil {
			if code == 0 { code = http.StatusInternalServerError }
			http.Error(w, err.Error(), code)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(data)
	})
	mux.HandleFunc("/api/rate/history", func(w http.ResponseWriter, r *http.Request) {
		if !checkAPIKey(w, r) || !rl.allow(r.RemoteAddr) {
			return
		}
		log.Printf("GET %s %s", r.URL.Path, r.URL.RawQuery)
		target := javaBaseURL + "/api/rate/history?" + r.URL.RawQuery
		data, code, err := simpleProxyGet(target, r.Header.Get("Authorization"))
		if err != nil {
			if code == 0 { code = http.StatusInternalServerError }
			http.Error(w, err.Error(), code)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(data)
	})
	mux.HandleFunc("/api/agent/register", func(w http.ResponseWriter, r *http.Request) {
		if !checkAPIKey(w, r) || !rl.allow(r.RemoteAddr) {
			return
		}
		log.Printf("POST %s", r.URL.Path)
		target := javaBaseURL + "/api/agent/register"
		body, _ := io.ReadAll(r.Body)
		data, code, err := httpPost(target, "", body)
		if err != nil {
			if code == 0 { code = http.StatusInternalServerError }
			http.Error(w, err.Error(), code)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(data)
	})
	mux.HandleFunc("/api/order/active", func(w http.ResponseWriter, r *http.Request) {
		if !checkAPIKey(w, r) || !rl.allow(r.RemoteAddr) {
			return
		}
		log.Printf("GET %s", r.URL.Path)
		target := javaBaseURL + "/api/order/active"
		data, code, err := simpleProxyGet(target, r.Header.Get("Authorization"))
		if err != nil {
			if code == 0 { code = http.StatusInternalServerError }
			http.Error(w, err.Error(), code)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(data)
	})
	mux.HandleFunc("/api/order/list", func(w http.ResponseWriter, r *http.Request) {
		if !checkAPIKey(w, r) || !rl.allow(r.RemoteAddr) {
			return
		}
		log.Printf("GET %s", r.URL.Path)
		target := javaBaseURL + "/api/order/list"
		data, code, err := simpleProxyGet(target, r.Header.Get("Authorization"))
		if err != nil {
			if code == 0 { code = http.StatusInternalServerError }
			http.Error(w, err.Error(), code)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(data)
	})
	mux.HandleFunc("/api/fund/balance", func(w http.ResponseWriter, r *http.Request) {
		if !checkAPIKey(w, r) || !rl.allow(r.RemoteAddr) {
			return
		}
		log.Printf("GET %s %s", r.URL.Path, r.URL.RawQuery)
		target := javaBaseURL + "/api/fund/balance?" + r.URL.RawQuery
		data, code, err := simpleProxyGet(target, r.Header.Get("Authorization"))
		if err != nil {
			if code == 0 { code = http.StatusInternalServerError }
			http.Error(w, err.Error(), code)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(data)
	})
	mux.HandleFunc("/api/order/create", func(w http.ResponseWriter, r *http.Request) {
		if !checkAPIKey(w, r) || !rl.allow(r.RemoteAddr) {
			return
		}
		// 获取 Authorization header 用于调试
		authHeader := r.Header.Get("Authorization")
		log.Printf("POST %s (Auth len: %d)", r.URL.Path, len(authHeader))
		
		target := javaBaseURL + "/api/order/create"
		body, _ := io.ReadAll(r.Body)
		data, code, err := httpPost(target, authHeader, body)
		if err != nil {
			if code == 0 { code = http.StatusInternalServerError }
			http.Error(w, err.Error(), code)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(data)
	})
	mux.HandleFunc("/api/order/buy", func(w http.ResponseWriter, r *http.Request) {
		if !checkAPIKey(w, r) || !rl.allow(r.RemoteAddr) {
			return
		}
		log.Printf("POST %s", r.URL.Path)
		target := javaBaseURL + "/api/order/buy"
		body, _ := io.ReadAll(r.Body)
		data, code, err := httpPost(target, r.Header.Get("Authorization"), body)
		if err != nil {
			if code == 0 { code = http.StatusInternalServerError }
			http.Error(w, err.Error(), code)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(data)
	})
	mux.HandleFunc("/api/order/sell", func(w http.ResponseWriter, r *http.Request) {
		if !checkAPIKey(w, r) || !rl.allow(r.RemoteAddr) {
			return
		}
		log.Printf("POST %s", r.URL.Path)
		target := javaBaseURL + "/api/order/sell"
		body, _ := io.ReadAll(r.Body)
		data, code, err := httpPost(target, r.Header.Get("Authorization"), body)
		if err != nil {
			if code == 0 { code = http.StatusInternalServerError }
			http.Error(w, err.Error(), code)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(data)
	})
	mux.HandleFunc("/api/order/cancel", func(w http.ResponseWriter, r *http.Request) {
		if !checkAPIKey(w, r) || !rl.allow(r.RemoteAddr) {
			return
		}
		log.Printf("POST %s", r.URL.Path)
		target := javaBaseURL + "/api/order/cancel"
		body, _ := io.ReadAll(r.Body)
		data, code, err := httpPost(target, r.Header.Get("Authorization"), body)
		if err != nil {
			if code == 0 { code = http.StatusInternalServerError }
			http.Error(w, err.Error(), code)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(data)
	})
	mux.HandleFunc("/api/agent/detail/", func(w http.ResponseWriter, r *http.Request) {
		if !checkAPIKey(w, r) || !rl.allow(r.RemoteAddr) {
			return
		}
		log.Printf("GET %s", r.URL.Path)
		target := javaBaseURL + r.URL.Path // 直接透传路径 /api/agent/detail/{agentId}
		data, code, err := simpleProxyGet(target, r.Header.Get("Authorization"))
		if err != nil {
			if code == 0 { code = http.StatusInternalServerError }
			http.Error(w, err.Error(), code)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(data)
	})
	mux.HandleFunc("/api/stream/rates", func(w http.ResponseWriter, r *http.Request) {
		if !checkAPIKey(w, r) || !rl.allow(r.RemoteAddr) {
			return
		}
		handleRatesSSE(w, r)
	})
	mux.HandleFunc("/api/stream/orders", func(w http.ResponseWriter, r *http.Request) {
		if !checkAPIKey(w, r) || !rl.allow(r.RemoteAddr) {
			return
		}
		handleOrdersSSE(w, r)
	})

	go func() {
		_ = http.ListenAndServe(":8090", mux)
	}()
	startRatesPoller(javaBaseURL)
	startOrdersPoller(javaBaseURL)
}

func simpleProxyGet(url, auth string) ([]byte, int, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 500, err
	}
	if auth != "" {
		req.Header.Set("Authorization", auth)
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 500, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return body, resp.StatusCode, nil
}

func checkAPIKey(w http.ResponseWriter, r *http.Request) bool {
	required := os.Getenv("MCP_API_KEY")
	if required == "" {
		return true
	}
	got := r.Header.Get("X-API-Key")
	if got != required {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return false
	}
	return true
}

type rateLimiter struct {
	limit     int
	interval  time.Duration
	counters  map[string]int
	resetTime time.Time
}

func newRateLimiter(limit int, interval time.Duration) *rateLimiter {
	return &rateLimiter{
		limit:     limit,
		interval:  interval,
		counters:  map[string]int{},
		resetTime: time.Now().Add(interval),
	}
}

func (rl *rateLimiter) allow(addr string) bool {
	now := time.Now()
	if now.After(rl.resetTime) {
		rl.counters = map[string]int{}
		rl.resetTime = now.Add(rl.interval)
	}
	host := addr
	if i := strings.LastIndex(addr, ":"); i > 0 {
		host = addr[:i]
	}
	rl.counters[host]++
	if rl.counters[host] > rl.limit {
		return false
	}
	return true
}

var ratesSubscribers = struct {
	clients map[chan []byte]struct{}
	last    []byte
}{clients: map[chan []byte]struct{}{}}

var ordersSubscribers = struct {
	clients map[chan []byte]struct{}
	last    []byte
}{clients: map[chan []byte]struct{}{}}

func handleRatesSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "stream unsupported", http.StatusInternalServerError)
		return
	}
	ch := make(chan []byte, 16)
	ratesSubscribers.clients[ch] = struct{}{}
	if ratesSubscribers.last != nil {
		fmt.Fprintf(w, "data: %s\n\n", string(ratesSubscribers.last))
		flusher.Flush()
	}
	notify := w.(http.CloseNotifier).CloseNotify()
	for {
		select {
		case data := <-ch:
			fmt.Fprintf(w, "data: %s\n\n", string(data))
			flusher.Flush()
		case <-notify:
			delete(ratesSubscribers.clients, ch)
			return
		}
	}
}

func handleOrdersSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "stream unsupported", http.StatusInternalServerError)
		return
	}
	ch := make(chan []byte, 16)
	ordersSubscribers.clients[ch] = struct{}{}
	if ordersSubscribers.last != nil {
		fmt.Fprintf(w, "data: %s\n\n", string(ordersSubscribers.last))
		flusher.Flush()
	}
	notify := w.(http.CloseNotifier).CloseNotify()
	for {
		select {
		case data := <-ch:
			fmt.Fprintf(w, "data: %s\n\n", string(data))
			flusher.Flush()
		case <-notify:
			delete(ordersSubscribers.clients, ch)
			return
		}
	}
}

func startRatesPoller(javaBaseURL string) {
	go func() {
		t := time.NewTicker(5 * time.Second)
		for range t.C {
			url := fmt.Sprintf("%s/rate/current?tradingPair=MON/USDC", javaBaseURL)
			data, err := httpGet(url, "")
			if err != nil {
				continue
			}
			b, _ := json.Marshal(data["data"])
			if !bytes.Equal(b, ratesSubscribers.last) {
				ratesSubscribers.last = b
				for ch := range ratesSubscribers.clients {
					select {
					case ch <- b:
					default:
					}
				}
			}
		}
	}()
}

func startOrdersPoller(javaBaseURL string) {
	go func() {
		t := time.NewTicker(5 * time.Second)
		for range t.C {
			url := fmt.Sprintf("%s/order/active", javaBaseURL)
			data, err := httpGet(url, "")
			if err != nil {
				continue
			}
			b, _ := json.Marshal(data["data"])
			if !bytes.Equal(b, ordersSubscribers.last) {
				ordersSubscribers.last = b
				for ch := range ordersSubscribers.clients {
					select {
					case ch <- b:
					default:
					}
				}
			}
		}
	}()
}
