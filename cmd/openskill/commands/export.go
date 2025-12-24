package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"openskill/pkg/skills"

	"github.com/spf13/cobra"
)

var exportFormat string
var exportOutput string

var ExportCmd = &cobra.Command{
	Use:   "export <skill-name>",
	Short: "Export a skill to different formats",
	Long: `Export a skill to JSON, YAML, or Markdown format.

Useful for sharing skills, backing up, or integrating with other tools.`,
	Args: cobra.ExactArgs(1),
	Example: `  openskill export code-review
  openskill export code-review --format json
  openskill export code-review --format yaml -o skill.yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		mgr := skills.NewManager()

		content, err := mgr.Export(name, exportFormat)
		if err != nil {
			return err
		}

		if exportOutput != "" {
			// Write to file
			dir := filepath.Dir(exportOutput)
			if dir != "." {
				if err := os.MkdirAll(dir, 0755); err != nil {
					return fmt.Errorf("failed to create directory: %w", err)
				}
			}
			if err := os.WriteFile(exportOutput, []byte(content), 0644); err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}
			fmt.Printf("✓ Exported '%s' to %s\n", name, exportOutput)
		} else {
			// Print to stdout
			fmt.Println(content)
		}

		return nil
	},
}

func init() {
	ExportCmd.Flags().StringVarP(&exportFormat, "format", "f", "yaml", "Export format (json, yaml, md)")
	ExportCmd.Flags().StringVarP(&exportOutput, "output", "o", "", "Output file path (default: stdout)")
}
