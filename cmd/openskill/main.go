package main

import (
	"fmt"
	"os"

	"openskill/cmd/openskill/commands"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "openskill",
	Short:   "Manage Claude skills",
	Long: `OpenSkill CLI - AI-Powered Skill Management for Claude

Create, manage, and share skills that enhance Claude's capabilities.
Skills are reusable behavior modules that define how Claude should
reason and act in specific domains.

Quick Start:
  openskill init              # Initialize in your project
  openskill template list     # Browse skill templates
  openskill add "my-skill"    # Create a new skill with AI
  openskill list              # View all skills`,
	Version: "0.3.0",
}

func init() {
	// Core commands
	rootCmd.AddCommand(commands.InitCmd)
	rootCmd.AddCommand(commands.AddCmd)
	rootCmd.AddCommand(commands.ListCmd)
	rootCmd.AddCommand(commands.ShowCmd)
	rootCmd.AddCommand(commands.EditCmd)
	rootCmd.AddCommand(commands.RemoveCmd)
	rootCmd.AddCommand(commands.ValidateCmd)
	rootCmd.AddCommand(commands.ConfigCmd)

	// Version history
	rootCmd.AddCommand(commands.HistoryCmd)
	rootCmd.AddCommand(commands.RollbackCmd)
	rootCmd.AddCommand(commands.DiffCmd)

	// Templates
	rootCmd.AddCommand(commands.TemplateCmd)

	// Import/Export
	rootCmd.AddCommand(commands.ExportCmd)
	rootCmd.AddCommand(commands.ImportCmd)

	// Testing
	rootCmd.AddCommand(commands.TestCmd)

	// AI-powered
	rootCmd.AddCommand(commands.ImproveCmd)
	rootCmd.AddCommand(commands.ExplainCmd)

	// Organization
	rootCmd.AddCommand(commands.TagCmd)
	rootCmd.AddCommand(commands.GroupCmd)
	rootCmd.AddCommand(commands.WorkspaceCmd)

	// Sync
	rootCmd.AddCommand(commands.SyncCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
