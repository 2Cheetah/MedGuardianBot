package groq

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateChatCompletion(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Error("Expected Authorization header with Bearer token")
		}

		// Mock response
		response := ChatCompletionResponse{
			ID:      "test-id",
			Object:  "chat.completion",
			Created: 1234567890,
			Model:   "mixtral-8x7b-32768",
			Choices: []struct {
				Message      Message `json:"message"`
				FinishReason string  `json:"finish_reason"`
			}{
				{
					Message: Message{
						Role:    "assistant",
						Content: "Test response",
					},
					FinishReason: "stop",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Error("couldn't encode to JSON")
		}
	}))
	defer mockServer.Close()

	client := NewClient("test-key")
	client.baseURL = mockServer.URL // Override base URL for testing

	req := ChatCompletionRequest{
		Model: "mixtral-8x7b-32768",
		Messages: []Message{
			{Role: "user", Content: "Hello"},
		},
	}

	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.Choices[0].Message.Content != "Test response" {
		t.Errorf("Expected 'Test response', got %s", resp.Choices[0].Message.Content)
	}
}
