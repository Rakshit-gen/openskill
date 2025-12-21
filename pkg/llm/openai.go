package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"openskill/pkg/config"
)

// OpenAIClient implements the Provider interface for OpenAI
type OpenAIClient struct {
	apiKey   string
	model    string
	endpoint string
}

// NewOpenAIClient creates a new OpenAI client
func NewOpenAIClient() *OpenAIClient {
	return &OpenAIClient{
		apiKey:   config.GetProviderAPIKey(string(ProviderOpenAI)),
		model:    config.GetProviderModel(string(ProviderOpenAI)),
		endpoint: ProviderEndpoints[ProviderOpenAI],
	}
}

func (c *OpenAIClient) Name() string {
	return "OpenAI"
}

func (c *OpenAIClient) IsConfigured() bool {
	return c.apiKey != ""
}

func (c *OpenAIClient) Generate(prompt string) (string, error) {
	req := chatRequest{
		Model: c.model,
		Messages: []message{
			{Role: "user", Content: prompt},
		},
	}

	body, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", c.endpoint, bytes.NewBuffer(body))
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("OpenAI API error: %s", resp.Status)
	}

	var result chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return result.Choices[0].Message.Content, nil
}
