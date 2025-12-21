package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"openskill/pkg/config"
)

// Client is the Groq API client (kept for backwards compatibility)
type Client struct {
	apiKey   string
	model    string
	endpoint string
}

// NewClient creates a new Groq client
func NewClient() *Client {
	return &Client{
		apiKey:   config.GetProviderAPIKey(string(ProviderGroq)),
		model:    config.GetProviderModel(string(ProviderGroq)),
		endpoint: ProviderEndpoints[ProviderGroq],
	}
}

func (c *Client) Name() string {
	return "Groq"
}

func (c *Client) IsConfigured() bool {
	return c.apiKey != ""
}

type chatRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (c *Client) Generate(prompt string) (string, error) {
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
		return "", fmt.Errorf("Groq API error: %s", resp.Status)
	}

	var result chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from Groq")
	}

	return result.Choices[0].Message.Content, nil
}
