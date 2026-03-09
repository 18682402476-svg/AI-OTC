package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// OpenAIStandardClient 通用 OpenAI 标准接口客户端
type OpenAIStandardClient struct {
	APIKey  string
	ModelID string
	URL     string
}

// NewOpenAIStandardClient 初始化客户端
func NewOpenAIStandardClient() *OpenAIStandardClient {
	url := os.Getenv("LLM_API_URL")
	if url == "" {
		url = os.Getenv("VOLC_API_URL")
	}
	if url == "" {
		url = "https://ark.cn-beijing.volces.com/api/v3/chat/completions"
	}

	apiKey := os.Getenv("LLM_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("VOLC_API_KEY")
	}

	modelID := os.Getenv("LLM_MODEL_ID")
	if modelID == "" {
		modelID = os.Getenv("VOLC_ENDPOINT_ID")
	}

	return &OpenAIStandardClient{
		APIKey:  apiKey,
		ModelID: modelID,
		URL:     url,
	}
}

// GetDecision 调用大语言模型 (LLM) 获取决策结果
func (c *OpenAIStandardClient) GetDecision(systemPrompt, userPrompt string) (*AgentAction, error) {
	if c.APIKey == "" || c.ModelID == "" {
		return nil, fmt.Errorf("环境变量 LLM_API_KEY/VOLC_API_KEY 或 LLM_MODEL_ID/VOLC_ENDPOINT_ID 未设置")
	}

	requestBody := map[string]interface{}{
		"model": c.ModelID,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userPrompt},
		},
		"response_format": map[string]string{"type": "json_object"},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	var apiResponse struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, err
	}

	if len(apiResponse.Choices) == 0 {
		return nil, fmt.Errorf("API 响应中没有 Choice 数据")
	}

	content := apiResponse.Choices[0].Message.Content
	
	// 清理 Markdown 代码块标记
	content = strings.TrimSpace(content)
	if strings.HasPrefix(content, "```json") {
		content = strings.TrimPrefix(content, "```json")
		content = strings.TrimSuffix(content, "```")
	} else if strings.HasPrefix(content, "```") {
		content = strings.TrimPrefix(content, "```")
		content = strings.TrimSuffix(content, "```")
	}
	content = strings.TrimSpace(content)

	var action AgentAction
	if err := json.Unmarshal([]byte(content), &action); err != nil {
		return nil, fmt.Errorf("解析 LLM 决策 JSON 失败: %v, 原文: %s", err, content)
	}

	return &action, nil
}
