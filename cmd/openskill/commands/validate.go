package commands

import (
	"fmt"
	"os"
	"strings"

	"openskill/pkg/skills"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var ValidateCmd = &cobra.Command{
	Use:   "validate <name>",
	Short: "Validate a skill's YAML structure and rules",
	Long: `Validate a skill file to ensure it has correct YAML syntax
and follows OpenSkill conventions.

Checks performed:
  • YAML syntax validation
  • Required fields (name, description)
  • Rule format and content
  • Best practices recommendations`,
	Args: cobra.ExactArgs(1),
	RunE: runValidate,
}

type ValidationResult struct {
	Errors   []string
	Warnings []string
}

func runValidate(cmd *cobra.Command, args []string) error {
	name := args[0]
	mgr := skills.NewManager()

	// Try to load the skill
	skill, err := mgr.Get(name)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("skill '%s' not found", name)
		}
		// YAML parsing error
		fmt.Printf("\n  ❌ Validation Failed: %s\n\n", name)
		fmt.Println("  YAML Syntax Error:")
		fmt.Printf("  └─ %v\n\n", err)
		return nil
	}

	result := validateSkill(skill.Name, skill.Description, skill.Rules)

	// Print results
	fmt.Println()
	if len(result.Errors) == 0 && len(result.Warnings) == 0 {
		fmt.Printf("  ✓ Skill '%s' is valid\n\n", name)
		printSkillSummary(skill.Name, skill.Description, skill.Rules)
		return nil
	}

	if len(result.Errors) > 0 {
		fmt.Printf("  ❌ Validation Failed: %s\n\n", name)
		fmt.Println("  Errors:")
		for _, e := range result.Errors {
			fmt.Printf("  └─ %s\n", e)
		}
		fmt.Println()
	}

	if len(result.Warnings) > 0 {
		if len(result.Errors) == 0 {
			fmt.Printf("  ⚠ Validation Passed with Warnings: %s\n\n", name)
		}
		fmt.Println("  Warnings:")
		for _, w := range result.Warnings {
			fmt.Printf("  └─ %s\n", w)
		}
		fmt.Println()
	}

	if len(result.Errors) == 0 {
		printSkillSummary(skill.Name, skill.Description, skill.Rules)
	}

	return nil
}

func validateSkill(name, description string, rules []string) ValidationResult {
	result := ValidationResult{}

	// Check required fields
	if name == "" {
		result.Errors = append(result.Errors, "Missing required field: name")
	} else {
		// Check name format
		if strings.Contains(name, " ") {
			result.Warnings = append(result.Warnings, "Skill name contains spaces - consider using hyphens (e.g., 'code-review')")
		}
		if len(name) > 50 {
			result.Warnings = append(result.Warnings, "Skill name is very long - consider a shorter, more memorable name")
		}
	}

	if description == "" {
		result.Errors = append(result.Errors, "Missing required field: description")
	} else {
		if len(description) < 10 {
			result.Warnings = append(result.Warnings, "Description is very short - add more detail for clarity")
		}
		if len(description) > 500 {
			result.Warnings = append(result.Warnings, "Description is very long - consider being more concise")
		}
	}

	// Check rules
	if len(rules) == 0 {
		result.Warnings = append(result.Warnings, "No rules defined - skills work better with specific behavioral rules")
	} else {
		for i, rule := range rules {
			if rule == "" {
				result.Errors = append(result.Errors, fmt.Sprintf("Rule %d is empty", i+1))
				continue
			}
			if len(rule) < 10 {
				result.Warnings = append(result.Warnings, fmt.Sprintf("Rule %d is very short - be more specific", i+1))
			}
			if len(rule) > 500 {
				result.Warnings = append(result.Warnings, fmt.Sprintf("Rule %d is very long - consider breaking into multiple rules", i+1))
			}
			// Check for common anti-patterns
			lowerRule := strings.ToLower(rule)
			if strings.HasPrefix(lowerRule, "be good") || strings.HasPrefix(lowerRule, "be nice") {
				result.Warnings = append(result.Warnings, fmt.Sprintf("Rule %d is vague - use specific, actionable instructions", i+1))
			}
		}

		if len(rules) > 20 {
			result.Warnings = append(result.Warnings, "Many rules defined - consider consolidating related rules")
		}
	}

	return result
}

func printSkillSummary(name, description string, rules []string) {
	fmt.Println("  Skill Summary:")
	fmt.Println("  ─────────────")
	fmt.Printf("  Name:        %s\n", name)
	fmt.Printf("  Description: %s\n", truncate(description, 60))
	fmt.Printf("  Rules:       %d defined\n", len(rules))
	fmt.Println()
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// ValidateYAML checks if a string is valid YAML
func ValidateYAML(content string) error {
	var data interface{}
	return yaml.Unmarshal([]byte(content), &data)
}
