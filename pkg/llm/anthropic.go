package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"openskill/pkg/config"
)

// AnthropicClient implements the Provider interface for Anthropic
type AnthropicClient struct {
	apiKey   string
	model    string
	endpoint string
}

// NewAnthropicClient creates a new Anthropic client
func NewAnthropicClient() *AnthropicClient {
	return &AnthropicClient{
		apiKey:   config.GetProviderAPIKey(string(ProviderAnthropic)),
		model:    config.GetProviderModel(string(ProviderAnthropic)),
		endpoint: ProviderEndpoints[ProviderAnthropic],
	}
}

func (c *AnthropicClient) Name() string {
	return "Anthropic"
}

func (c *AnthropicClient) IsConfigured() bool {
	return c.apiKey != ""
}

// Anthropic-specific request/response types
type anthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	Messages  []anthropicMessage `json:"messages"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (c *AnthropicClient) Generate(prompt string) (string, error) {
	req := anthropicRequest{
		Model:     c.model,
		MaxTokens: 4096,
		Messages: []anthropicMessage{
			{Role: "user", Content: prompt},
		},
	}

	body, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", c.endpoint, bytes.NewBuffer(body))
	httpReq.Header.Set("x-api-key", c.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result anthropicResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		if result.Error != nil {
			return "", fmt.Errorf("Anthropic API error: %s", result.Error.Message)
		}
		return "", fmt.Errorf("Anthropic API error: %s", resp.Status)
	}

	if len(result.Content) == 0 {
		return "", fmt.Errorf("no response from Anthropic")
	}

	return result.Content[0].Text, nil
}
