package commands

import (
	"fmt"

	"openskill/pkg/skills"

	"github.com/spf13/cobra"
)

var ShowCmd = &cobra.Command{
	Use:   "show <name>",
	Short: "Show skill details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		mgr := skills.NewManager()
		skill, err := mgr.Get(name)
		if err != nil {
			return fmt.Errorf("skill '%s' not found", name)
		}

		fmt.Printf("Name: %s\n", skill.Name)
		fmt.Printf("Description: %s\n", skill.Description)
		if len(skill.Rules) > 0 {
			fmt.Println("Rules:")
			for i, r := range skill.Rules {
				fmt.Printf("  %d. %s\n", i+1, r)
			}
		}
		return nil
	},
}
