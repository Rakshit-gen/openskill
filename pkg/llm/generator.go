package llm

import (
	"encoding/json"
	"fmt"
	"strings"

	"openskill/pkg/core"
)

type Generator struct {
	client *Client
}

func NewGenerator() *Generator {
	return &Generator{client: NewClient()}
}

func (g *Generator) IsAvailable() bool {
	return g.client.IsConfigured()
}

func (g *Generator) EnhanceSkill(name, description string) (*core.Skill, error) {
	prompt := fmt.Sprintf(`Create a Claude skill definition for: "%s"
User description: %s

Generate a JSON response with:
1. A clear, detailed description
2. multiple in depth specific rules/guidelines for Claude to follow
3. Make sure to include all the details of the user description in the rules
4. add general rules for the skill that are not specific to the user description

Response format (JSON only, no markdown):
{
  "description": "...",
  "rules": ["rule1", "rule2", "rule3", ...]
}`, name, description)

	response, err := g.client.Generate(prompt)
	if err != nil {
		return nil, err
	}

	// Clean response - remove markdown code blocks if present
	response = strings.TrimSpace(response)
	response = strings.TrimPrefix(response, "```json")
	response = strings.TrimPrefix(response, "```")
	response = strings.TrimSuffix(response, "```")
	response = strings.TrimSpace(response)

	var result struct {
		Description string   `json:"description"`
		Rules       []string `json:"rules"`
	}

	if err := json.Unmarshal([]byte(response), &result); err != nil {
		return nil, fmt.Errorf("failed to parse LLM response: %w", err)
	}

	return &core.Skill{
		Name:        name,
		Description: result.Description,
		Rules:       result.Rules,
	}, nil
}
