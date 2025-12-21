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
	prompt := fmt.Sprintf(`You are an expert prompt engineer creating a comprehensive Claude skill definition.

Skill Name: "%s"
User's Intent: %s

Create a FLAWLESS and EXHAUSTIVE skill definition that will make Claude an absolute expert in this domain.

Requirements:
1. DESCRIPTION: Write a rich, detailed description (3-4 sentences) that:
   - Clearly defines the skill's purpose and scope
   - Explains what makes this skill valuable
   - Sets clear expectations for Claude's behavior

2. RULES: Generate 8-12 comprehensive rules that cover:
   - Core principles and best practices for this skill
   - Specific actionable guidelines Claude must follow
   - Edge cases and how to handle them
   - Quality standards and success criteria
   - Common pitfalls to avoid
   - User experience considerations
   - Output format and communication style
   - When to ask for clarification vs make assumptions

Each rule should be:
- Specific and actionable (not vague)
- Self-contained and clear
- Focused on one concept
- Written as a directive ("Always...", "Never...", "When X, do Y...")

Response format (JSON only, no markdown, no code blocks):
{
  "description": "...",
  "rules": ["rule1", "rule2", "rule3", "rule4", "rule5", "rule6", "rule7", "rule8", ...]
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
