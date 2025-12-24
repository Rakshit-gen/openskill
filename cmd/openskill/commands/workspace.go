package commands

import (
	"fmt"
	"strings"

	"openskill/pkg/core"
	"openskill/pkg/skills"

	"github.com/spf13/cobra"
)

var WorkspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Manage project-specific skill configuration",
	Long: `Configure which skills are active for the current project.

Workspaces allow you to:
- Enable/disable specific skills per project
- Override skill variables for project-specific needs
- Group related projects with similar skill needs`,
}

var workspaceInitCmd = &cobra.Command{
	Use:   "init [name]",
	Short: "Initialize a workspace in the current project",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := "default"
		if len(args) > 0 {
			name = args[0]
		}

		// Check if workspace already exists
		existing, _ := skills.LoadWorkspace()
		if existing != nil {
			return fmt.Errorf("workspace already exists: %s", existing.Name)
		}

		workspace := &core.Workspace{
			Name:        name,
			Description: fmt.Sprintf("Workspace for %s", name),
			Skills:      []string{},
			Groups:      []string{},
			Overrides:   make(map[string]map[string]string),
		}

		if err := skills.SaveWorkspace(workspace); err != nil {
			return fmt.Errorf("failed to create workspace: %w", err)
		}

		fmt.Printf("✓ Workspace '%s' created\n", name)
		fmt.Println("  Location: .claude/workspace.yaml")
		fmt.Println("\n  Use 'openskill workspace add <skill>' to add skills")

		return nil
	},
}

var workspaceShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current workspace configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		workspace, err := skills.LoadWorkspace()
		if err != nil {
			return err
		}
		if workspace == nil {
			fmt.Println("No workspace configured.")
			fmt.Println("Use 'openskill workspace init' to create one.")
			return nil
		}

		fmt.Printf("\nWorkspace: %s\n", workspace.Name)
		fmt.Println("═══════════════════════════════════════════════════")

		if workspace.Description != "" {
			fmt.Printf("Description: %s\n\n", workspace.Description)
		}

		if len(workspace.Skills) > 0 {
			fmt.Println("Enabled Skills:")
			for _, skill := range workspace.Skills {
				fmt.Printf("  • %s\n", skill)
			}
			fmt.Println()
		} else {
			fmt.Println("No skills enabled.\n")
		}

		if len(workspace.Groups) > 0 {
			fmt.Println("Enabled Groups:")
			for _, group := range workspace.Groups {
				fmt.Printf("  • %s\n", group)
			}
			fmt.Println()
		}

		if len(workspace.Overrides) > 0 {
			fmt.Println("Variable Overrides:")
			for skill, vars := range workspace.Overrides {
				fmt.Printf("  %s:\n", skill)
				for k, v := range vars {
					fmt.Printf("    %s = %s\n", k, v)
				}
			}
			fmt.Println()
		}

		return nil
	},
}

var workspaceAddCmd = &cobra.Command{
	Use:   "add <skill-name>",
	Short: "Add a skill to the workspace",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		skillName := args[0]

		workspace, err := skills.LoadWorkspace()
		if err != nil {
			return err
		}
		if workspace == nil {
			return fmt.Errorf("no workspace configured. Use 'openskill workspace init'")
		}

		// Check if skill exists
		mgr := skills.NewManager()
		if _, err := mgr.Get(skillName); err != nil {
			return fmt.Errorf("skill '%s' not found", skillName)
		}

		// Check if already added
		for _, s := range workspace.Skills {
			if strings.EqualFold(s, skillName) {
				return fmt.Errorf("skill '%s' is already in the workspace", skillName)
			}
		}

		workspace.Skills = append(workspace.Skills, skillName)
		if err := skills.SaveWorkspace(workspace); err != nil {
			return err
		}

		fmt.Printf("✓ Added '%s' to workspace\n", skillName)
		return nil
	},
}

var workspaceRemoveCmd = &cobra.Command{
	Use:   "remove <skill-name>",
	Short: "Remove a skill from the workspace",
	Aliases: []string{"rm"},
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		skillName := args[0]

		workspace, err := skills.LoadWorkspace()
		if err != nil {
			return err
		}
		if workspace == nil {
			return fmt.Errorf("no workspace configured")
		}

		// Find and remove skill
		found := false
		var newSkills []string
		for _, s := range workspace.Skills {
			if strings.EqualFold(s, skillName) {
				found = true
			} else {
				newSkills = append(newSkills, s)
			}
		}

		if !found {
			return fmt.Errorf("skill '%s' not in workspace", skillName)
		}

		workspace.Skills = newSkills
		if err := skills.SaveWorkspace(workspace); err != nil {
			return err
		}

		fmt.Printf("✓ Removed '%s' from workspace\n", skillName)
		return nil
	},
}

var workspaceSetCmd = &cobra.Command{
	Use:   "set <skill-name> <variable> <value>",
	Short: "Set a variable override for a skill",
	Args:  cobra.ExactArgs(3),
	Example: `  openskill workspace set code-review max_issues 10
  openskill workspace set testing framework jest`,
	RunE: func(cmd *cobra.Command, args []string) error {
		skillName := args[0]
		variable := args[1]
		value := args[2]

		workspace, err := skills.LoadWorkspace()
		if err != nil {
			return err
		}
		if workspace == nil {
			return fmt.Errorf("no workspace configured")
		}

		if workspace.Overrides == nil {
			workspace.Overrides = make(map[string]map[string]string)
		}
		if workspace.Overrides[skillName] == nil {
			workspace.Overrides[skillName] = make(map[string]string)
		}

		workspace.Overrides[skillName][variable] = value
		if err := skills.SaveWorkspace(workspace); err != nil {
			return err
		}

		fmt.Printf("✓ Set %s.%s = %s\n", skillName, variable, value)
		return nil
	},
}

func init() {
	WorkspaceCmd.AddCommand(workspaceInitCmd)
	WorkspaceCmd.AddCommand(workspaceShowCmd)
	WorkspaceCmd.AddCommand(workspaceAddCmd)
	WorkspaceCmd.AddCommand(workspaceRemoveCmd)
	WorkspaceCmd.AddCommand(workspaceSetCmd)
}
