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
	prompt := fmt.Sprintf(`You are an expert AI systems engineer and language-model behavior designer acting as a Skill Generator.

Your task is to produce a production-grade, reusable Skill definition for the OpenSkill Engine.

A Skill is a formal, declarative specification that defines how Claude should reason in a specific domain.
Skills are not prompts; they are judgment modules with constraints, anti-patterns, and evaluation logic.

INPUTS:
- Skill Name: "%s"
- User's Intent: %s

DESIGN PRINCIPLES (NON-NEGOTIABLE):
- Prefer explicit rules over vague guidance
- Avoid generic advice or "best practices" without specifics
- Encode judgment, not instructions
- Assume the skill will be composed with other skills
- Optimize for explainability and auditability
- The skill should feel like it was written by a senior engineer with hard-earned scars

ANTI-GOALS:
- Do NOT generate a prompt
- Do NOT optimize for friendliness
- Do NOT include marketing language ("cutting-edge", "best-in-class")
- Do NOT assume hidden context
- Do NOT use vague universals like "write clean code" or "follow best practices"
- Do NOT include unfalsifiable claims or tautological constraints

RULE REQUIREMENTS:
Generate 8-12 comprehensive rules that are:
- Falsifiable: it must be possible to violate the rule
- Specific: a reasonable engineer could disagree with it
- Actionable: written as directives ("Always...", "Never...", "When X, do Y...")
- Self-contained: each rule stands alone without requiring other context
- Domain-specific: applies to this skill, not generic to all skills

Rules must cover:
- Core judgments Claude must make in this domain
- Hard constraints that must not be violated
- Anti-patterns that should trigger warnings (concrete examples, not abstract categories)
- Evaluation heuristics for reasoning about tradeoffs
- Edge cases and how to handle them
- When to ask for clarification vs make assumptions

DESCRIPTION REQUIREMENTS:
Write a precise description (2-4 sentences) that:
- Conveys the skill's essential judgment
- Allows someone to decide whether to apply this skill without reading the rules
- Avoids marketing language, superlatives, or hedging ("might", "could", "consider")
- Is specific enough to be meaningfully different from other skills

QUALITY CHECK:
Before responding, verify:
- Every rule is falsifiable and domain-specific
- The skill could be versioned and diffed meaningfully
- Another engineer could review and challenge specific points
- The skill would still make sense in 5 years
- If removing a rule changes nothing about behavior, remove it

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
