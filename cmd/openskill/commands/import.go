package commands

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"openskill/pkg/skills"

	"github.com/spf13/cobra"
)

var importFormat string
var importName string

var ImportCmd = &cobra.Command{
	Use:   "import <file-or-url>",
	Short: "Import a skill from a file or URL",
	Long: `Import a skill from a JSON or YAML file, or from a URL.

Supports importing from:
- Local files (.json, .yaml, .yml)
- GitHub Gists
- Raw URLs
- Stdin (use - as filename)`,
	Args: cobra.ExactArgs(1),
	Example: `  openskill import skill.json
  openskill import skill.yaml --name my-skill
  openskill import https://gist.githubusercontent.com/.../skill.json
  cat skill.json | openskill import -`,
	RunE: func(cmd *cobra.Command, args []string) error {
		source := args[0]
		var content string
		var detectedFormat string

		if source == "-" {
			// Read from stdin
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("failed to read stdin: %w", err)
			}
			content = string(data)
			detectedFormat = importFormat
		} else if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
			// Fetch from URL
			resp, err := http.Get(source)
			if err != nil {
				return fmt.Errorf("failed to fetch URL: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("HTTP error: %s", resp.Status)
			}

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("failed to read response: %w", err)
			}
			content = string(data)

			// Detect format from URL
			if strings.HasSuffix(source, ".json") {
				detectedFormat = "json"
			} else if strings.HasSuffix(source, ".yaml") || strings.HasSuffix(source, ".yml") {
				detectedFormat = "yaml"
			} else {
				detectedFormat = importFormat
			}
		} else {
			// Read from file
			data, err := os.ReadFile(source)
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}
			content = string(data)

			// Detect format from extension
			ext := strings.ToLower(filepath.Ext(source))
			switch ext {
			case ".json":
				detectedFormat = "json"
			case ".yaml", ".yml":
				detectedFormat = "yaml"
			default:
				detectedFormat = importFormat
			}
		}

		if importFormat != "" {
			detectedFormat = importFormat
		}

		mgr := skills.NewManager()
		skill, err := mgr.Import(content, detectedFormat)
		if err != nil {
			return err
		}

		// Override name if specified
		if importName != "" {
			skill.Name = importName
		}

		if skill.Name == "" {
			return fmt.Errorf("skill name is required (use --name flag)")
		}

		if err := mgr.Add(skill); err != nil {
			return err
		}

		fmt.Printf("✓ Imported skill: %s\n", skill.Name)
		fmt.Printf("  Description: %s\n", skill.Description)
		fmt.Printf("  Rules: %d\n", len(skill.Rules))

		return nil
	},
}

func init() {
	ImportCmd.Flags().StringVarP(&importFormat, "format", "f", "yaml", "Import format (json, yaml)")
	ImportCmd.Flags().StringVarP(&importName, "name", "n", "", "Override skill name")
}
