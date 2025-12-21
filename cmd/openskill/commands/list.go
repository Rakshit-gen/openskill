package commands

import (
	"fmt"

	"openskill/pkg/skills"

	"github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all skills",
	RunE: func(cmd *cobra.Command, args []string) error {
		mgr := skills.NewManager()
		skillList, err := mgr.List()
		if err != nil {
			return err
		}

		if len(skillList) == 0 {
			fmt.Println("No skills found. Add one with: openskill add <name> -d \"description\"")
			return nil
		}

		fmt.Printf("Skills (%d):\n\n", len(skillList))
		for _, s := range skillList {
			fmt.Printf("  %s\n", s.Name)
			fmt.Printf("    %s\n", s.Description)
			if len(s.Rules) > 0 {
				fmt.Printf("    Rules: %d\n", len(s.Rules))
			}
			fmt.Println()
		}
		return nil
	},
}
