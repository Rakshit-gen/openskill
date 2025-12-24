package commands

import (
	"fmt"
	"strings"

	"openskill/pkg/core"
	"openskill/pkg/skills"

	"github.com/spf13/cobra"
)

var listTag string
var listGroup string
var listVerbose bool

var ListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all skills",
	Long: `List all skills, optionally filtered by tag or group.

Use --tag to filter by tag, --group to filter by group.
Use --verbose to see more details about each skill.`,
	Example: `  openskill list
  openskill list --tag security
  openskill list --group development
  openskill list -v`,
	RunE: func(cmd *cobra.Command, args []string) error {
		mgr := skills.NewManager()

		var skillList []core.Skill
		var err error

		if listTag != "" {
			skillList, err = mgr.ListByTag(listTag)
			if err != nil {
				return err
			}
		} else if listGroup != "" {
			skillList, err = mgr.ListByGroup(listGroup)
			if err != nil {
				return err
			}
		} else {
			skillList, err = mgr.List()
			if err != nil {
				return err
			}
		}

		if len(skillList) == 0 {
			if listTag != "" {
				fmt.Printf("No skills found with tag '%s'\n", listTag)
			} else if listGroup != "" {
				fmt.Printf("No skills found in group '%s'\n", listGroup)
			} else {
				fmt.Println("No skills found. Add one with: openskill add <name> -d \"description\"")
			}
			return nil
		}

		// Print header
		header := "Skills"
		if listTag != "" {
			header = fmt.Sprintf("Skills tagged '%s'", listTag)
		} else if listGroup != "" {
			header = fmt.Sprintf("Skills in group '%s'", listGroup)
		}
		fmt.Printf("\n%s (%d):\n", header, len(skillList))
		fmt.Println("─────────────────────────────────────────────────────")

		for _, s := range skillList {
			fmt.Printf("\n  %s", s.Name)

			// Show version and template info inline
			if s.Version != "" {
				fmt.Printf(" (v%s)", s.Version)
			}
			if s.Template != "" {
				fmt.Printf(" [from: %s]", s.Template)
			}
			fmt.Println()

			// Description
			desc := s.Description
			if !listVerbose && len(desc) > 70 {
				desc = desc[:67] + "..."
			}
			fmt.Printf("    %s\n", desc)

			// Metadata line
			var meta []string
			if len(s.Rules) > 0 {
				meta = append(meta, fmt.Sprintf("%d rules", len(s.Rules)))
			}
			if s.Group != "" && listGroup == "" {
				meta = append(meta, fmt.Sprintf("group: %s", s.Group))
			}
			if len(s.Tags) > 0 && listTag == "" {
				meta = append(meta, fmt.Sprintf("tags: %s", strings.Join(s.Tags, ", ")))
			}
			if len(meta) > 0 {
				fmt.Printf("    [%s]\n", strings.Join(meta, " | "))
			}

			// Verbose mode: show more details
			if listVerbose {
				if s.Author != "" {
					fmt.Printf("    Author: %s\n", s.Author)
				}
				if s.Extends != "" {
					fmt.Printf("    Extends: %s\n", s.Extends)
				}
				if len(s.Includes) > 0 {
					fmt.Printf("    Includes: %s\n", strings.Join(s.Includes, ", "))
				}
				if len(s.Chain) > 0 {
					fmt.Printf("    Chain: %s\n", strings.Join(s.Chain, " → "))
				}
				if s.OutputFormat != "" {
					fmt.Printf("    Output: %s\n", s.OutputFormat)
				}
				if s.Context != nil {
					contextParts := []string{}
					if len(s.Context.Files) > 0 {
						contextParts = append(contextParts, fmt.Sprintf("%d files", len(s.Context.Files)))
					}
					if len(s.Context.Globs) > 0 {
						contextParts = append(contextParts, fmt.Sprintf("%d globs", len(s.Context.Globs)))
					}
					if len(s.Context.Commands) > 0 {
						contextParts = append(contextParts, fmt.Sprintf("%d commands", len(s.Context.Commands)))
					}
					if len(contextParts) > 0 {
						fmt.Printf("    Context: %s\n", strings.Join(contextParts, ", "))
					}
				}
				if s.Hooks != nil {
					hookParts := []string{}
					if len(s.Hooks.Pre) > 0 {
						hookParts = append(hookParts, fmt.Sprintf("%d pre", len(s.Hooks.Pre)))
					}
					if len(s.Hooks.Post) > 0 {
						hookParts = append(hookParts, fmt.Sprintf("%d post", len(s.Hooks.Post)))
					}
					if len(hookParts) > 0 {
						fmt.Printf("    Hooks: %s\n", strings.Join(hookParts, ", "))
					}
				}
			}
		}
		fmt.Println()
		return nil
	},
}

func init() {
	ListCmd.Flags().StringVarP(&listTag, "tag", "t", "", "Filter by tag")
	ListCmd.Flags().StringVarP(&listGroup, "group", "g", "", "Filter by group")
	ListCmd.Flags().BoolVarP(&listVerbose, "verbose", "v", false, "Show detailed information")
}
