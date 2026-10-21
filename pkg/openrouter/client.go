package openrouter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type LlmClient interface {
	GenerateAnswer(ctx context.Context, systemPrompt, userPrompt string) (string, error)
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float32   `json:"temperature,omitempty"`
}

type ResponseError struct {
	Message string `json:"message"`
	Code    any    `json:"code"`
}

func (obj *ResponseError) Error() string {
	if obj.Code != nil {
		return fmt.Sprintf("code %v: %s", obj.Code, obj.Message)
	}
	return obj.Message
}

func (obj *ResponseError) UnmarshalJSON(data []byte) error {
	var errStr string
	if err := json.Unmarshal(data, &errStr); err == nil {
		obj.Message = errStr
		return nil
	}
	type rawError struct {
		Message string `json:"message"`
		Code    any    `json:"code"`
	}
	var raw rawError
	if err := json.Unmarshal(data, &raw); err == nil {
		obj.Message = raw.Message
		obj.Code = raw.Code
		if obj.Message == "" && len(data) > 0 {
			obj.Message = string(data)
		}
		return nil
	}
	obj.Message = string(data)
	return nil
}

type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
	Error *ResponseError `json:"error,omitempty"`
}

type Client struct {
	apiKey  string
	baseUrl string
	model   string
	client  *http.Client
}

func NewClient(apiKey string, baseUrl string, model string, client *http.Client) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("api key is required")
	}
	if baseUrl == "" {
		return nil, fmt.Errorf("base url is required")
	}
	if model == "" {
		return nil, fmt.Errorf("model is required")
	}
	return &Client{
		apiKey:  apiKey,
		baseUrl: baseUrl,
		model:   model,
		client:  client,
	}, nil
}

func (obj *Client) GenerateAnswer(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	requestConfig := ChatRequest{
		Model: obj.model,
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: 0.3,
	}
	data, err := json.Marshal(requestConfig)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, obj.baseUrl+"/chat/completions", bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+obj.apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("HTTP-Referer", "https://github.com/Timur1414/Smart-Catch-Up")
	req.Header.Set("X-Title", "VK-Smart-Catch-Up")

	resp, err := obj.client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errResp struct {
			Error   *ResponseError `json:"error"`
			Message string         `json:"message"`
		}
		if err := json.Unmarshal(body, &errResp); err == nil {
			if errResp.Error != nil {
				return "", fmt.Errorf("openrouter error (HTTP %d): %s", resp.StatusCode, errResp.Error.Error())
			}
			if errResp.Message != "" {
				return "", fmt.Errorf("openrouter error (HTTP %d): %s", resp.StatusCode, errResp.Message)
			}
		}
		return "", fmt.Errorf("openrouter error (HTTP %d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	var response ChatResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}
	if response.Error != nil {
		return "", fmt.Errorf(response.Error.Message)
	}
	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from open router")
	}
	return response.Choices[0].Message.Content, nil
}
