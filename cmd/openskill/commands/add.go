package commands

import (
	"fmt"

	"openskill/pkg/core"
	"openskill/pkg/llm"
	"openskill/pkg/skills"

	"github.com/spf13/cobra"
)

var (
	addDesc   string
	addRules  []string
	addManual bool
)

var AddCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Add a new skill (uses AI to generate content)",
	Args:  cobra.ExactArgs(1),
	Example: `  openskill add "code-review" -d "Reviews code"
  openskill add "bug-finder" -d "Finds bugs" --manual -r "Check nulls"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		if addDesc == "" {
			return fmt.Errorf("description is required (-d flag)")
		}

		var skill *core.Skill

		// Use LLM to enhance the skill unless --manual is set
		if !addManual {
			gen := llm.NewGenerator()
			if !gen.IsAvailable() {
				return fmt.Errorf("API key not configured. Set it with:\n\n  openskill config set api-key\n\nOr use --manual flag to skip AI generation")
			}
			fmt.Printf("Generating skill with %s...\n", gen.ProviderName())
			enhanced, err := gen.EnhanceSkill(name, addDesc)
			if err != nil {
				return fmt.Errorf("AI generation failed: %w", err)
			}
			skill = enhanced
		} else {
			// Manual mode
			skill = &core.Skill{
				Name:        name,
				Description: addDesc,
				Rules:       addRules,
			}
		}

		mgr := skills.NewManager()
		if err := mgr.Add(skill); err != nil {
			return err
		}

		fmt.Printf("\n✓ Added skill: %s\n", skill.Name)
		fmt.Printf("  Description: %s\n", skill.Description)
		if len(skill.Rules) > 0 {
			fmt.Println("  Rules:")
			for i, r := range skill.Rules {
				fmt.Printf("    %d. %s\n", i+1, r)
			}
		}
		return nil
	},
}

func init() {
	AddCmd.Flags().StringVarP(&addDesc, "desc", "d", "", "Skill description (required)")
	AddCmd.Flags().StringArrayVarP(&addRules, "rule", "r", nil, "Add a rule (manual mode only)")
	AddCmd.Flags().BoolVar(&addManual, "manual", false, "Skip AI generation, use provided values")
}
