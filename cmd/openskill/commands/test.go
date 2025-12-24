package commands

import (
	"fmt"
	"strings"

	"openskill/pkg/llm"
	"openskill/pkg/skills"

	"github.com/spf13/cobra"
)

var testPrompt string
var testMock bool

var TestCmd = &cobra.Command{
	Use:   "test <skill-name>",
	Short: "Test a skill with a sample prompt",
	Long: `Test a skill by running it against a sample prompt.

This helps validate that a skill works as expected before using it in production.
Use --mock to see how the skill would be applied without making an API call.`,
	Args: cobra.ExactArgs(1),
	Example: `  openskill test code-review --prompt "Review this function: func add(a, b int) int { return a + b }"
  openskill test commit-message --mock`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		mgr := skills.NewManager()

		skill, err := mgr.Get(name)
		if err != nil {
			return fmt.Errorf("skill '%s' not found", name)
		}

		fmt.Printf("\nTesting skill: %s\n", skill.Name)
		fmt.Println("═══════════════════════════════════════════════════")

		// Build the skill context
		var context strings.Builder
		context.WriteString(fmt.Sprintf("You are operating with the '%s' skill.\n\n", skill.Name))
		context.WriteString(fmt.Sprintf("Description: %s\n\n", skill.Description))

		if len(skill.Rules) > 0 {
			context.WriteString("Rules you must follow:\n")
			for i, rule := range skill.Rules {
				context.WriteString(fmt.Sprintf("%d. %s\n", i+1, rule))
			}
		}

		if testMock {
			fmt.Println("\n[Mock Mode - No API call made]")
			fmt.Println("\nSkill context that would be sent:")
			fmt.Println("───────────────────────────────────")
			fmt.Println(context.String())

			if testPrompt != "" {
				fmt.Println("\nUser prompt:")
				fmt.Println("───────────────────────────────────")
				fmt.Println(testPrompt)
			}
			fmt.Println()
			return nil
		}

		if testPrompt == "" {
			return fmt.Errorf("--prompt is required (or use --mock for dry run)")
		}

		// Make actual API call
		gen := llm.NewGenerator()
		if !gen.IsAvailable() {
			return fmt.Errorf("no AI provider configured. Use 'openskill config set api-key' or --mock flag")
		}

		fullPrompt := context.String() + "\n\nUser request:\n" + testPrompt

		fmt.Printf("\nRunning with %s...\n", gen.ProviderName())
		fmt.Println("───────────────────────────────────")

		// Use the provider directly for custom prompt
		response, err := gen.Provider().Generate(fullPrompt)
		if err != nil {
			return fmt.Errorf("API call failed: %w", err)
		}

		fmt.Println("\nResponse:")
		fmt.Println("───────────────────────────────────")
		fmt.Println(response)
		fmt.Println()

		return nil
	},
}

func init() {
	TestCmd.Flags().StringVarP(&testPrompt, "prompt", "p", "", "Test prompt to run against the skill")
	TestCmd.Flags().BoolVar(&testMock, "mock", false, "Mock mode - show skill context without API call")
}
