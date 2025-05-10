package groq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

const (
	defaultBaseURL = "https://api.groq.com/"
	chatEndpoint   = "openai/v1/chat/completions"
)

const (
	RoleSystem    = "system"    // System messages help set behavior
	RoleUser      = "user"      // User messages are the inputs from the end user
	RoleAssistant = "assistant" // Assistant messages are the model's responses
	RoleFunction  = "function"  // Function messages represent function call results
)

type Client struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewClient(apiKey string) *Client {
	slog.Info("creating groq client")
	return &Client{
		apiKey:  apiKey,
		baseURL: defaultBaseURL,
		client:  &http.Client{},
	}
}

type Message struct {
	Role    string `json:"role"` // Must be one of: "system", "user", "assistant", "function"
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
}

type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Message      Message `json:"message"`
		FinishReason string  `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func (c *Client) CreateChatCompletion(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error) {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", defaultBaseURL+chatEndpoint, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result ChatCompletionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &result, nil
}

func (c *Client) GetCompletion(prompt string, systemPrompt string) (string, error) {
	slog.Info("entering ParseSchedule func")
	req := ChatCompletionRequest{
		Model: "gemma2-9b-it",
		Messages: []Message{
			{
				Role:    RoleUser,
				Content: prompt,
			},
			{
				Role:    RoleSystem,
				Content: systemPrompt,
			},
		},
	}
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", defaultBaseURL+chatEndpoint, bytes.NewReader(jsonBody))
	slog.Info("about to HTTP call", "httpReq", httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	slog.Info("received http response from API", "body", body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result ChatCompletionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result.Choices[0].Message.Content, nil
}

func (c *Client) ParseTextToISO8601(text string) (string, error) {
	systemPrompt := "Extract data from the input prompt and convert it to `ISO 8601` date format. No explanation, no additional data, just a text of date in a form of `YYYY-MM-DD`, for example `2025-10-25`."
	resp, err := c.GetCompletion(text, systemPrompt)
	if err != nil {
		return "", fmt.Errorf("couldn't get chat completion from prompt: %s, error: %w", text, err)
	}
	res := strings.TrimSpace(resp)
	return res, nil
}
