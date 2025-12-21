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
	Version: "0.2.0",
}

func init() {
	rootCmd.AddCommand(commands.InitCmd)
	rootCmd.AddCommand(commands.AddCmd)
	rootCmd.AddCommand(commands.ListCmd)
	rootCmd.AddCommand(commands.ShowCmd)
	rootCmd.AddCommand(commands.EditCmd)
	rootCmd.AddCommand(commands.RemoveCmd)
	rootCmd.AddCommand(commands.ValidateCmd)
	rootCmd.AddCommand(commands.HistoryCmd)
	rootCmd.AddCommand(commands.RollbackCmd)
	rootCmd.AddCommand(commands.ConfigCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
