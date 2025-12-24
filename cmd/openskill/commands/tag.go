package commands

import (
	"fmt"
	"strings"

	"openskill/pkg/skills"

	"github.com/spf13/cobra"
)

var TagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Manage skill tags",
	Long: `Organize skills with tags for flexible categorization.

Tags allow you to:
- Categorize skills by multiple dimensions
- Filter skills by tag
- Discover related skills`,
}

var tagListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tags",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		mgr := skills.NewManager()
		tags, err := mgr.GetAllTags()
		if err != nil {
			return err
		}

		if len(tags) == 0 {
			fmt.Println("No tags defined.")
			fmt.Println("Add tags to a skill with: openskill tag add <skill> <tag>")
			return nil
		}

		fmt.Println("\nAll Tags:")
		fmt.Println("─────────────────────────────────────")

		for _, tag := range tags {
			skillsWithTag, _ := mgr.ListByTag(tag)
			fmt.Printf("  %-20s (%d skills)\n", tag, len(skillsWithTag))
		}
		fmt.Println()

		return nil
	},
}

var tagShowCmd = &cobra.Command{
	Use:   "show <tag-name>",
	Short: "Show skills with a tag",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tagName := args[0]
		mgr := skills.NewManager()

		skillsWithTag, err := mgr.ListByTag(tagName)
		if err != nil {
			return err
		}

		if len(skillsWithTag) == 0 {
			return fmt.Errorf("no skills found with tag '%s'", tagName)
		}

		fmt.Printf("\nSkills tagged '%s':\n", tagName)
		fmt.Println("═══════════════════════════════════════════════════")

		for _, skill := range skillsWithTag {
			fmt.Printf("\n  %s\n", skill.Name)
			fmt.Printf("    %s\n", truncateText(skill.Description, 60))
			if len(skill.Tags) > 1 {
				var otherTags []string
				for _, t := range skill.Tags {
					if !strings.EqualFold(t, tagName) {
						otherTags = append(otherTags, t)
					}
				}
				if len(otherTags) > 0 {
					fmt.Printf("    Other tags: %s\n", strings.Join(otherTags, ", "))
				}
			}
		}
		fmt.Println()

		return nil
	},
}

var tagAddCmd = &cobra.Command{
	Use:   "add <skill-name> <tag>...",
	Short: "Add tags to a skill",
	Args:  cobra.MinimumNArgs(2),
	Example: `  openskill tag add code-review quality security
  openskill tag add api-design rest architecture`,
	RunE: func(cmd *cobra.Command, args []string) error {
		skillName := args[0]
		newTags := args[1:]

		mgr := skills.NewManager()
		skill, err := mgr.Get(skillName)
		if err != nil {
			return fmt.Errorf("skill '%s' not found", skillName)
		}

		// Add new tags (avoiding duplicates)
		existingTags := make(map[string]bool)
		for _, t := range skill.Tags {
			existingTags[strings.ToLower(t)] = true
		}

		addedTags := []string{}
		for _, tag := range newTags {
			if !existingTags[strings.ToLower(tag)] {
				skill.Tags = append(skill.Tags, tag)
				addedTags = append(addedTags, tag)
				existingTags[strings.ToLower(tag)] = true
			}
		}

		if len(addedTags) == 0 {
			fmt.Println("All tags already exist on this skill.")
			return nil
		}

		if err := mgr.Edit(skillName, skill); err != nil {
			return err
		}

		fmt.Printf("✓ Added tags to '%s': %s\n", skillName, strings.Join(addedTags, ", "))
		return nil
	},
}

var tagRemoveCmd = &cobra.Command{
	Use:   "remove <skill-name> <tag>...",
	Short: "Remove tags from a skill",
	Aliases: []string{"rm"},
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		skillName := args[0]
		tagsToRemove := args[1:]

		mgr := skills.NewManager()
		skill, err := mgr.Get(skillName)
		if err != nil {
			return fmt.Errorf("skill '%s' not found", skillName)
		}

		// Remove specified tags
		removeSet := make(map[string]bool)
		for _, t := range tagsToRemove {
			removeSet[strings.ToLower(t)] = true
		}

		var newTags []string
		removedTags := []string{}
		for _, t := range skill.Tags {
			if removeSet[strings.ToLower(t)] {
				removedTags = append(removedTags, t)
			} else {
				newTags = append(newTags, t)
			}
		}

		if len(removedTags) == 0 {
			fmt.Println("None of the specified tags exist on this skill.")
			return nil
		}

		skill.Tags = newTags
		if err := mgr.Edit(skillName, skill); err != nil {
			return err
		}

		fmt.Printf("✓ Removed tags from '%s': %s\n", skillName, strings.Join(removedTags, ", "))
		return nil
	},
}

func init() {
	TagCmd.AddCommand(tagListCmd)
	TagCmd.AddCommand(tagShowCmd)
	TagCmd.AddCommand(tagAddCmd)
	TagCmd.AddCommand(tagRemoveCmd)
}
