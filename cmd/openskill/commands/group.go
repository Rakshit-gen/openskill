package commands

import (
	"fmt"
	"strings"

	"openskill/pkg/skills"

	"github.com/spf13/cobra"
)

var GroupCmd = &cobra.Command{
	Use:   "group",
	Short: "Manage skill groups",
	Long: `Organize skills into groups for easier management.

Groups allow you to:
- Bundle related skills together
- Enable/disable multiple skills at once
- Organize skills by project or domain`,
}

var groupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all groups",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		mgr := skills.NewManager()
		groups, err := mgr.GetAllGroups()
		if err != nil {
			return err
		}

		if len(groups) == 0 {
			fmt.Println("No groups defined.")
			fmt.Println("Add a group to a skill with: openskill edit <skill> --group <name>")
			return nil
		}

		fmt.Println("\nSkill Groups:")
		fmt.Println("─────────────────────────────────────")

		for _, group := range groups {
			skillsInGroup, _ := mgr.ListByGroup(group)
			fmt.Printf("\n  %s (%d skills)\n", group, len(skillsInGroup))
			for _, skill := range skillsInGroup {
				fmt.Printf("    • %s\n", skill.Name)
			}
		}
		fmt.Println()

		return nil
	},
}

var groupShowCmd = &cobra.Command{
	Use:   "show <group-name>",
	Short: "Show skills in a group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		groupName := args[0]
		mgr := skills.NewManager()

		skillsInGroup, err := mgr.ListByGroup(groupName)
		if err != nil {
			return err
		}

		if len(skillsInGroup) == 0 {
			return fmt.Errorf("group '%s' not found or empty", groupName)
		}

		fmt.Printf("\nGroup: %s\n", groupName)
		fmt.Println("═══════════════════════════════════════════════════")
		fmt.Printf("Skills: %d\n\n", len(skillsInGroup))

		for _, skill := range skillsInGroup {
			fmt.Printf("  %s\n", skill.Name)
			fmt.Printf("    %s\n", truncateText(skill.Description, 60))
			fmt.Printf("    Rules: %d", len(skill.Rules))
			if len(skill.Tags) > 0 {
				fmt.Printf("  Tags: %s", strings.Join(skill.Tags, ", "))
			}
			fmt.Println("\n")
		}

		return nil
	},
}

var groupSetCmd = &cobra.Command{
	Use:   "set <skill-name> <group-name>",
	Short: "Add a skill to a group",
	Args:  cobra.ExactArgs(2),
	Example: `  openskill group set code-review development
  openskill group set security-audit security`,
	RunE: func(cmd *cobra.Command, args []string) error {
		skillName := args[0]
		groupName := args[1]

		mgr := skills.NewManager()
		skill, err := mgr.Get(skillName)
		if err != nil {
			return fmt.Errorf("skill '%s' not found", skillName)
		}

		// Save version before modifying
		if err := mgr.SaveVersion(skillName); err != nil {
			fmt.Printf("Warning: Could not save version: %v\n", err)
		}

		skill.Group = groupName
		if err := mgr.Edit(skillName, skill); err != nil {
			return err
		}

		fmt.Printf("✓ Added '%s' to group '%s'\n", skillName, groupName)
		return nil
	},
}

var groupUnsetCmd = &cobra.Command{
	Use:   "unset <skill-name>",
	Short: "Remove a skill from its group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		skillName := args[0]

		mgr := skills.NewManager()
		skill, err := mgr.Get(skillName)
		if err != nil {
			return fmt.Errorf("skill '%s' not found", skillName)
		}

		if skill.Group == "" {
			return fmt.Errorf("skill '%s' is not in any group", skillName)
		}

		oldGroup := skill.Group
		skill.Group = ""
		if err := mgr.Edit(skillName, skill); err != nil {
			return err
		}

		fmt.Printf("✓ Removed '%s' from group '%s'\n", skillName, oldGroup)
		return nil
	},
}

func truncateText(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func init() {
	GroupCmd.AddCommand(groupListCmd)
	GroupCmd.AddCommand(groupShowCmd)
	GroupCmd.AddCommand(groupSetCmd)
	GroupCmd.AddCommand(groupUnsetCmd)
}
