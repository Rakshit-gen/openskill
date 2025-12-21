package commands

import (
	"fmt"

	"openskill/pkg/skills"

	"github.com/spf13/cobra"
)

var (
	editDesc  string
	editRules []string
	editName  string
)

var EditCmd = &cobra.Command{
	Use:   "edit <name>",
	Short: "Edit an existing skill",
	Args:  cobra.ExactArgs(1),
	Example: `  openskill edit "code-review" -d "New description"
  openskill edit "bug-finder" -r "New rule 1" -r "New rule 2"
  openskill edit "old-name" --name "new-name"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		mgr := skills.NewManager()
		skill, err := mgr.Get(name)
		if err != nil {
			return fmt.Errorf("skill '%s' not found", name)
		}

		// Update fields if provided
		if editName != "" {
			skill.Name = editName
		}
		if editDesc != "" {
			skill.Description = editDesc
		}
		if len(editRules) > 0 {
			skill.Rules = editRules
		}

		if err := mgr.Edit(name, skill); err != nil {
			return err
		}

		fmt.Printf("✓ Updated skill: %s\n", skill.Name)
		return nil
	},
}

func init() {
	EditCmd.Flags().StringVar(&editName, "name", "", "New name for the skill")
	EditCmd.Flags().StringVarP(&editDesc, "desc", "d", "", "New description")
	EditCmd.Flags().StringArrayVarP(&editRules, "rule", "r", nil, "Replace rules (can be used multiple times)")
}
