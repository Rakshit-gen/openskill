package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"openskill/pkg/llm"
	"openskill/pkg/skills"

	"github.com/spf13/cobra"
)

var improveApply bool

var ImproveCmd = &cobra.Command{
	Use:   "improve <skill-name>",
	Short: "Use AI to suggest improvements to a skill",
	Long: `Analyze a skill and suggest improvements using AI.

This command reviews your skill's rules and description, identifying:
- Vague or non-actionable rules
- Missing edge cases
- Potential conflicts between rules
- Opportunities for better specificity`,
	Args: cobra.ExactArgs(1),
	Example: `  openskill improve code-review
  openskill improve code-review --apply`,
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

		fmt.Printf("Analyzing skill '%s' with %s...\n\n", name, gen.ProviderName())

		// Build analysis prompt
		var rulesText strings.Builder
		for i, rule := range skill.Rules {
			rulesText.WriteString(fmt.Sprintf("%d. %s\n", i+1, rule))
		}

		prompt := fmt.Sprintf(`Analyze this skill definition and suggest improvements.

Skill Name: %s
Description: %s

Current Rules:
%s

Analyze the skill and provide:
1. Overall assessment (1-2 sentences)
2. Specific issues with existing rules (if any)
3. Suggested new or improved rules
4. Any missing edge cases or considerations

Return your response as JSON:
{
  "assessment": "Overall assessment here",
  "issues": ["issue 1", "issue 2"],
  "improved_rules": ["improved rule 1", "improved rule 2", ...],
  "improved_description": "Better description if needed, or empty string"
}`, skill.Name, skill.Description, rulesText.String())

		response, err := gen.Provider().Generate(prompt)
		if err != nil {
			return fmt.Errorf("AI analysis failed: %w", err)
		}

		// Clean and parse response
		response = strings.TrimSpace(response)
		response = strings.TrimPrefix(response, "```json")
		response = strings.TrimPrefix(response, "```")
		response = strings.TrimSuffix(response, "```")
		response = strings.TrimSpace(response)

		var result struct {
			Assessment          string   `json:"assessment"`
			Issues              []string `json:"issues"`
			ImprovedRules       []string `json:"improved_rules"`
			ImprovedDescription string   `json:"improved_description"`
		}

		if err := json.Unmarshal([]byte(response), &result); err != nil {
			fmt.Println("AI Response:")
			fmt.Println(response)
			return nil
		}

		// Display results
		fmt.Println("Assessment:")
		fmt.Println("───────────────────────────────────")
		fmt.Printf("  %s\n\n", result.Assessment)

		if len(result.Issues) > 0 {
			fmt.Println("Issues Found:")
			fmt.Println("───────────────────────────────────")
			for _, issue := range result.Issues {
				fmt.Printf("  • %s\n", issue)
			}
			fmt.Println()
		}

		if len(result.ImprovedRules) > 0 {
			fmt.Println("Suggested Rules:")
			fmt.Println("───────────────────────────────────")
			for i, rule := range result.ImprovedRules {
				fmt.Printf("  %d. %s\n", i+1, rule)
			}
			fmt.Println()
		}

		if result.ImprovedDescription != "" && result.ImprovedDescription != skill.Description {
			fmt.Println("Suggested Description:")
			fmt.Println("───────────────────────────────────")
			fmt.Printf("  %s\n\n", result.ImprovedDescription)
		}

		if improveApply && len(result.ImprovedRules) > 0 {
			// Save version before applying
			if err := mgr.SaveVersion(name); err != nil {
				fmt.Printf("Warning: Could not save version history: %v\n", err)
			}

			// Apply improvements
			skill.Rules = result.ImprovedRules
			if result.ImprovedDescription != "" {
				skill.Description = result.ImprovedDescription
			}

			if err := mgr.Edit(name, skill); err != nil {
				return fmt.Errorf("failed to apply improvements: %w", err)
			}

			fmt.Println("✓ Improvements applied!")
			fmt.Println("  Use 'openskill rollback' to revert if needed.")
		} else if len(result.ImprovedRules) > 0 {
			fmt.Println("Run with --apply to apply these improvements.")
		}

		return nil
	},
}

func init() {
	ImproveCmd.Flags().BoolVar(&improveApply, "apply", false, "Apply the suggested improvements")
}
