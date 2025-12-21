package commands

import (
	"fmt"

	"openskill/pkg/skills"

	"github.com/spf13/cobra"
)

var RemoveCmd = &cobra.Command{
	Use:     "remove <name>",
	Aliases: []string{"rm"},
	Short:   "Remove a skill",
	Args:    cobra.ExactArgs(1),
	Example: `  openskill remove "code-review"
  openskill rm "bug-finder"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		mgr := skills.NewManager()
		if err := mgr.Remove(name); err != nil {
			return err
		}

		fmt.Printf("✓ Removed skill: %s\n", name)
		return nil
	},
}
