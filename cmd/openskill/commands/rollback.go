package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var RollbackCmd = &cobra.Command{
	Use:   "rollback <name> <version>",
	Short: "Restore a skill to a previous version",
	Long: `Restore a skill to a previous version from its history.

The current version will be saved to history before restoring.
Use 'openskill history <name>' to see available versions.`,
	Args: cobra.ExactArgs(2),
	RunE: runRollback,
}

func runRollback(cmd *cobra.Command, args []string) error {
	name := args[0]
	versionStr := args[1]

	// Parse version number
	version, err := strconv.Atoi(strings.TrimPrefix(versionStr, "v"))
	if err != nil {
		return fmt.Errorf("invalid version: %s (expected a number like '1' or 'v1')", versionStr)
	}

	safeName := strings.ReplaceAll(strings.ToLower(name), " ", "-")
	skillDir := filepath.Join(".claude/skills", safeName)
	skillPath := filepath.Join(skillDir, "SKILL.md")
	historyPath := filepath.Join(historyDir, safeName)
	versionFile := filepath.Join(historyPath, fmt.Sprintf("SKILL.v%d.md", version))

	// Check if skill exists
	if _, err := os.Stat(skillPath); os.IsNotExist(err) {
		return fmt.Errorf("skill '%s' not found", name)
	}

	// Check if version exists
	if _, err := os.Stat(versionFile); os.IsNotExist(err) {
		return fmt.Errorf("version %d not found for skill '%s'", version, name)
	}

	// Save current version to history before rollback
	if err := SaveVersion(name); err != nil {
		fmt.Printf("Warning: failed to backup current version: %v\n", err)
	}

	// Read the old version
	content, err := os.ReadFile(versionFile)
	if err != nil {
		return fmt.Errorf("failed to read version %d: %w", version, err)
	}

	// Restore it
	if err := os.WriteFile(skillPath, content, 0644); err != nil {
		return fmt.Errorf("failed to restore version: %w", err)
	}

	fmt.Println()
	fmt.Printf("  ✓ Restored '%s' to version %d\n", name, version)
	fmt.Println()
	fmt.Println("  The previous version has been saved to history.")
	fmt.Printf("  Use 'openskill show %s' to view the restored skill.\n", name)
	fmt.Println()

	return nil
}
