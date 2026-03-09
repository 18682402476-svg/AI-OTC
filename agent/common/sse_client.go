package common

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func StartRatesListener(tradingPair string) {
	base := os.Getenv("MCP_SERVER_URL")
	if base == "" {
		base = "http://localhost:8090/api"
	}
	url := fmt.Sprintf("%s/stream/rates?tradingPair=%s", base, tradingPair)
	go sseConsume(url, "")
}

func StartOrdersListener(token string) {
	base := os.Getenv("MCP_SERVER_URL")
	if base == "" {
		base = "http://localhost:8090/api"
	}
	url := fmt.Sprintf("%s/stream/orders", base)
	go sseConsume(url, token)
}

func sseConsume(url, token string) {
	for {
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Accept", "text/event-stream")
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		if k := os.Getenv("MCP_API_KEY"); k != "" {
			req.Header.Set("X-API-Key", k)
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}
		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			if strings.HasPrefix(line, "data: ") {
				data := strings.TrimSpace(strings.TrimPrefix(line, "data: "))
				log.Printf("SSE %s -> %s", url, data)
			}
		}
		resp.Body.Close()
		time.Sleep(2 * time.Second)
	}
}
