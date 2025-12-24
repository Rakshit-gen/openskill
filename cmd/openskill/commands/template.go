package commands

import (
	"fmt"
	"strings"

	"openskill/pkg/core"
	"openskill/pkg/skills"

	"github.com/spf13/cobra"
)

var TemplateCmd = &cobra.Command{
	Use:   "template",
	Short: "Manage skill templates",
	Long:  `List, show, and use built-in skill templates to quickly create new skills.`,
}

var templateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available templates",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		templates := skills.GetBuiltinTemplates()

		if len(templates) == 0 {
			fmt.Println("No templates available")
			return nil
		}

		// Group by category
		categories := make(map[string][]core.SkillTemplate)
		for _, t := range templates {
			categories[t.Category] = append(categories[t.Category], t)
		}

		fmt.Println("\nAvailable Templates:")
		fmt.Println("────────────────────")

		for category, temps := range categories {
			fmt.Printf("\n  %s:\n", strings.ToUpper(category))
			for _, t := range temps {
				fmt.Printf("    %-20s %s\n", t.Name, t.Description)
			}
		}
		fmt.Println()

		return nil
	},
}

var templateShowCmd = &cobra.Command{
	Use:   "show <name>",
	Short: "Show template details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		templates := skills.GetBuiltinTemplates()

		var found *core.SkillTemplate
		for _, t := range templates {
			if strings.EqualFold(t.Name, name) {
				found = &t
				break
			}
		}

		if found == nil {
			return fmt.Errorf("template '%s' not found", name)
		}

		fmt.Printf("\nTemplate: %s\n", found.Name)
		fmt.Println("─────────────────────────────────────")
		fmt.Printf("Category:    %s\n", found.Category)
		fmt.Printf("Description: %s\n\n", found.Description)
		fmt.Printf("Skill Description:\n  %s\n\n", found.Skill.Description)

		if len(found.Skill.Tags) > 0 {
			fmt.Printf("Tags: %s\n\n", strings.Join(found.Skill.Tags, ", "))
		}

		fmt.Println("Rules:")
		for i, rule := range found.Skill.Rules {
			fmt.Printf("  %d. %s\n", i+1, rule)
		}
		fmt.Println()

		return nil
	},
}

var templateUseCmd = &cobra.Command{
	Use:   "use <template> [skill-name]",
	Short: "Create a skill from a template",
	Args:  cobra.RangeArgs(1, 2),
	Example: `  openskill template use code-review
  openskill template use commit-message my-commit-helper`,
	RunE: func(cmd *cobra.Command, args []string) error {
		templateName := args[0]
		templates := skills.GetBuiltinTemplates()

		var found *core.SkillTemplate
		for _, t := range templates {
			if strings.EqualFold(t.Name, templateName) {
				found = &t
				break
			}
		}

		if found == nil {
			return fmt.Errorf("template '%s' not found", templateName)
		}

		// Create a copy of the skill
		skill := found.Skill
		skill.Template = found.Name

		// Use custom name if provided
		if len(args) > 1 {
			skill.Name = args[1]
		}

		mgr := skills.NewManager()
		if err := mgr.Add(&skill); err != nil {
			return err
		}

		fmt.Printf("\n✓ Created skill '%s' from template '%s'\n", skill.Name, found.Name)
		fmt.Printf("  Location: .claude/skills/%s/SKILL.md\n\n", strings.ToLower(skill.Name))

		return nil
	},
}

func init() {
	TemplateCmd.AddCommand(templateListCmd)
	TemplateCmd.AddCommand(templateShowCmd)
	TemplateCmd.AddCommand(templateUseCmd)
}
