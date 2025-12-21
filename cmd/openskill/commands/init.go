package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"openskill/pkg/config"

	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize OpenSkill in your project",
	Long: `Initialize OpenSkill by setting up the skills directory and configuration.

This command will:
  1. Create the .claude/skills/ directory for storing skill definitions
  2. Optionally set up your Groq API key for AI-powered skill generation
  3. Create an example skill to get you started

After running this command, Claude will automatically discover and use
any skills you add to the .claude/skills/ directory.`,
	RunE: runInit,
}

func runInit(cmd *cobra.Command, args []string) error {
	fmt.Println()
	fmt.Println("  ╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("  ║                                                           ║")
	fmt.Println("  ║   ⚡ OpenSkill - Claude Skill Manager                     ║")
	fmt.Println("  ║                                                           ║")
	fmt.Println("  ╚═══════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Step 1: Create skills directory
	fmt.Println("  [1/3] Setting up skills directory...")
	skillsDir := ".claude/skills"
	if err := os.MkdirAll(skillsDir, 0755); err != nil {
		return fmt.Errorf("failed to create skills directory: %w", err)
	}
	fmt.Printf("        ✓ Created %s/\n", skillsDir)
	fmt.Println()

	// Step 2: Check/Set API key
	fmt.Println("  [2/3] Checking API configuration...")
	apiKey := config.GetAPIKey()
	if apiKey == "" {
		fmt.Println("        No Groq API key found.")
		fmt.Println()
		fmt.Print("        Would you like to set up your API key now? (y/n): ")

		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response == "y" || response == "yes" {
			fmt.Print("        Enter your Groq API key: ")
			key, _ := reader.ReadString('\n')
			key = strings.TrimSpace(key)

			if key != "" {
				cfg, err := config.Load()
				if err != nil {
					cfg = &config.Config{}
				}
				cfg.GroqAPIKey = key
				if err := config.Save(cfg); err != nil {
					fmt.Printf("        ⚠ Failed to save API key: %v\n", err)
				} else {
					fmt.Println("        ✓ API key saved to ~/.openskill/config.yaml")
				}
			}
		} else {
			fmt.Println("        ⚠ Skipped. You can set it later with: openskill config set api-key")
		}
	} else {
		masked := apiKey[:4] + "..." + apiKey[len(apiKey)-4:]
		fmt.Printf("        ✓ API key configured (%s)\n", masked)
	}
	fmt.Println()

	// Step 3: Create example skill (directory-based with SKILL.md)
	fmt.Println("  [3/3] Creating example skill...")
	exampleDir := skillsDir + "/example"
	examplePath := exampleDir + "/SKILL.md"
	if _, err := os.Stat(exampleDir); os.IsNotExist(err) {
		// Create skill directory
		if err := os.MkdirAll(exampleDir, 0755); err != nil {
			fmt.Printf("        ⚠ Failed to create example skill directory: %v\n", err)
		} else {
			exampleContent := `---
name: example
description: An example skill to demonstrate the OpenSkill format
---

# example

An example skill to demonstrate the OpenSkill format.

## Rules

- Be helpful and concise in all responses
- Provide code examples when they would clarify the explanation
- Explain your reasoning step by step when solving problems
- Ask clarifying questions when the request is ambiguous
`
			if err := os.WriteFile(examplePath, []byte(exampleContent), 0644); err != nil {
				fmt.Printf("        ⚠ Failed to create example skill: %v\n", err)
			} else {
				fmt.Println("        ✓ Created example/SKILL.md")
			}
		}
	} else {
		fmt.Println("        ✓ Example skill already exists")
	}
	fmt.Println()

	// Print usage guide
	fmt.Println("  ╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("  ║                    Setup Complete!                        ║")
	fmt.Println("  ╚═══════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("  How Claude discovers your skills:")
	fmt.Println("  ─────────────────────────────────")
	fmt.Println("  Claude reads SKILL.md files from .claude/skills/<skill-name>/")
	fmt.Println("  when you start a conversation in this directory.")
	fmt.Println()
	fmt.Println("  Quick Start:")
	fmt.Println("  ─────────────")
	fmt.Println("  • Add a skill:       openskill add \"code review\"")
	fmt.Println("  • List skills:       openskill list")
	fmt.Println("  • Show a skill:      openskill show code-review")
	fmt.Println("  • Edit a skill:      openskill edit code-review")
	fmt.Println("  • Validate a skill:  openskill validate code-review")
	fmt.Println("  • Remove a skill:    openskill remove code-review")
	fmt.Println()
	fmt.Println("  Skill File Format (.claude/skills/<name>/SKILL.md):")
	fmt.Println("  ──────────────────────────────────────────────────")
	fmt.Println("  ---")
	fmt.Println("  name: skill-name")
	fmt.Println("  description: What this skill does")
	fmt.Println("  ---")
	fmt.Println()
	fmt.Println("  # skill-name")
	fmt.Println()
	fmt.Println("  ## Rules")
	fmt.Println("  - Rule 1: Be specific about behavior")
	fmt.Println("  - Rule 2: Define constraints and requirements")
	fmt.Println()
	fmt.Println("  Need help? Run: openskill --help")
	fmt.Println()

	return nil
}
