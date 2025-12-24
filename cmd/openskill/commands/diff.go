package commands

import (
	"fmt"
	"strings"

	"openskill/pkg/skills"

	"github.com/spf13/cobra"
)

var diffVersion1 int
var diffVersion2 int

var DiffCmd = &cobra.Command{
	Use:   "diff <skill-name>",
	Short: "Show differences between skill versions",
	Long: `Compare different versions of a skill to see what changed.

By default, compares the current version with the most recent saved version.
Use --v1 and --v2 flags to compare specific versions.`,
	Args: cobra.ExactArgs(1),
	Example: `  openskill diff code-review
  openskill diff code-review --v1 1 --v2 2
  openskill diff code-review --v1 3`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		mgr := skills.NewManager()

		// Check skill exists
		if _, err := mgr.Get(name); err != nil {
			return fmt.Errorf("skill '%s' not found", name)
		}

		// Get versions
		versions, err := mgr.GetVersions(name)
		if err != nil {
			return err
		}

		if len(versions) == 0 && diffVersion1 == 0 && diffVersion2 == 0 {
			fmt.Println("No version history available for this skill.")
			fmt.Println("Use 'openskill edit' to create versions.")
			return nil
		}

		// Default: compare current (0) with latest saved version
		v1 := diffVersion1
		v2 := diffVersion2
		if v1 == 0 && v2 == 0 && len(versions) > 0 {
			v1 = versions[0].Version
			v2 = 0 // current
		}

		content1, content2, err := mgr.Diff(name, v1, v2)
		if err != nil {
			return err
		}

		// Display diff
		label1 := "current"
		if v1 != 0 {
			label1 = fmt.Sprintf("v%d", v1)
		}
		label2 := "current"
		if v2 != 0 {
			label2 = fmt.Sprintf("v%d", v2)
		}

		fmt.Printf("\nComparing %s (%s) with %s (%s)\n", name, label1, name, label2)
		fmt.Println("═══════════════════════════════════════════════════")

		// Simple line-by-line diff
		lines1 := strings.Split(content1, "\n")
		lines2 := strings.Split(content2, "\n")

		// Find differences
		maxLines := len(lines1)
		if len(lines2) > maxLines {
			maxLines = len(lines2)
		}

		hasDiff := false
		for i := 0; i < maxLines; i++ {
			var l1, l2 string
			if i < len(lines1) {
				l1 = lines1[i]
			}
			if i < len(lines2) {
				l2 = lines2[i]
			}

			if l1 != l2 {
				hasDiff = true
				if l1 != "" && (i >= len(lines2) || l1 != l2) {
					fmt.Printf("- %s\n", l1)
				}
				if l2 != "" && (i >= len(lines1) || l1 != l2) {
					fmt.Printf("+ %s\n", l2)
				}
			}
		}

		if !hasDiff {
			fmt.Println("\nNo differences found.")
		}

		fmt.Println()
		return nil
	},
}

func init() {
	DiffCmd.Flags().IntVar(&diffVersion1, "v1", 0, "First version to compare (0 = current)")
	DiffCmd.Flags().IntVar(&diffVersion2, "v2", 0, "Second version to compare (0 = current)")
}
