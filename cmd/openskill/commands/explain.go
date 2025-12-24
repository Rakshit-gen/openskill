package commands

import (
	"fmt"
	"strings"

	"openskill/pkg/llm"
	"openskill/pkg/skills"

	"github.com/spf13/cobra"
)

var explainVerbose bool

var ExplainCmd = &cobra.Command{
	Use:   "explain <skill-name>",
	Short: "Get an AI-powered explanation of what a skill does",
	Long: `Use AI to generate a human-readable explanation of what a skill does.

This is helpful for understanding complex skills or onboarding new team members.`,
	Args: cobra.ExactArgs(1),
	Example: `  openskill explain code-review
  openskill explain security-audit --verbose`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		mgr := skills.NewManager()

		skill, err := mgr.Get(name)
		if err != nil {
			return fmt.Errorf("skill '%s' not found", name)
		}

		gen := llm.NewGenerator()
		if !gen.IsAvailable() {
			return fmt.Errorf("no AI provider configured. Use 'openskill config set api-key'")
		}

		fmt.Printf("Explaining skill '%s' with %s...\n\n", name, gen.ProviderName())

		// Build explanation prompt
		var rulesText strings.Builder
		for i, rule := range skill.Rules {
			rulesText.WriteString(fmt.Sprintf("%d. %s\n", i+1, rule))
		}

		verboseInstructions := ""
		if explainVerbose {
			verboseInstructions = `
Also include:
- Example scenarios where this skill would be applied
- Potential edge cases the skill handles
- How this skill might interact with other skills`
		}

		prompt := fmt.Sprintf(`Explain this skill in plain language for a developer who hasn't seen it before.

Skill Name: %s
Description: %s

Rules:
%s

Write a clear, concise explanation that covers:
1. What this skill is designed to do (1-2 sentences)
2. Key behaviors it enforces
3. What makes it effective
%s

Use simple language and avoid jargon. Format the response with clear sections.`,
			skill.Name, skill.Description, rulesText.String(), verboseInstructions)

		response, err := gen.Provider().Generate(prompt)
		if err != nil {
			return fmt.Errorf("AI explanation failed: %w", err)
		}

		fmt.Printf("Skill: %s\n", skill.Name)
		fmt.Println("═══════════════════════════════════════════════════")
		fmt.Println(strings.TrimSpace(response))
		fmt.Println()

		return nil
	},
}

func init() {
	ExplainCmd.Flags().BoolVarP(&explainVerbose, "verbose", "v", false, "Include examples and edge cases")
}
