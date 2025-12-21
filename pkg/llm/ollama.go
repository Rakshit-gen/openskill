package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"openskill/pkg/config"
)

// OllamaClient implements the Provider interface for Ollama (local)
type OllamaClient struct {
	model    string
	endpoint string
}

// NewOllamaClient creates a new Ollama client
func NewOllamaClient() *OllamaClient {
	endpoint := config.GetOllamaEndpoint()
	if endpoint == "" {
		endpoint = ProviderEndpoints[ProviderOllama]
	}
	return &OllamaClient{
		model:    config.GetProviderModel(string(ProviderOllama)),
		endpoint: endpoint,
	}
}

func (c *OllamaClient) Name() string {
	return "Ollama"
}

func (c *OllamaClient) IsConfigured() bool {
	// Ollama doesn't need an API key, just check if endpoint is reachable
	// For now, always return true and let it fail on actual request
	return true
}

// Ollama-specific request/response types
type ollamaRequest struct {
	Model    string          `json:"model"`
	Messages []ollamaMessage `json:"messages"`
	Stream   bool            `json:"stream"`
}

type ollamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaResponse struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
	Error string `json:"error,omitempty"`
}

func (c *OllamaClient) Generate(prompt string) (string, error) {
	req := ollamaRequest{
		Model: c.model,
		Messages: []ollamaMessage{
			{Role: "user", Content: prompt},
		},
		Stream: false,
	}

	body, _ := json.Marshal(req)
	httpReq, err := http.NewRequest("POST", c.endpoint, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("Ollama not reachable at %s: %w", c.endpoint, err)
	}
	defer resp.Body.Close()

	var result ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		if result.Error != "" {
			return "", fmt.Errorf("Ollama error: %s", result.Error)
		}
		return "", fmt.Errorf("Ollama API error: %s", resp.Status)
	}

	if result.Message.Content == "" {
		return "", fmt.Errorf("no response from Ollama")
	}

	return result.Message.Content, nil
}
