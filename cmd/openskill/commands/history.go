package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const historyDir = ".claude/skills/.history"

var HistoryCmd = &cobra.Command{
	Use:   "history <name>",
	Short: "Show version history of a skill",
	Long: `Display the version history of a skill, showing all previous versions
that have been saved. Each edit creates a new version that can be restored.

Use 'openskill rollback <name> <version>' to restore a previous version.`,
	Args: cobra.ExactArgs(1),
	RunE: runHistory,
}

type VersionInfo struct {
	Version   int
	Timestamp time.Time
	Filename  string
}

func runHistory(cmd *cobra.Command, args []string) error {
	name := args[0]
	safeName := strings.ReplaceAll(strings.ToLower(name), " ", "-")

	// Check if skill exists (directory-based structure)
	skillDir := filepath.Join(".claude/skills", safeName)
	skillPath := filepath.Join(skillDir, "SKILL.md")
	if _, err := os.Stat(skillPath); os.IsNotExist(err) {
		return fmt.Errorf("skill '%s' not found", name)
	}

	// Get skill file info
	info, err := os.Stat(skillPath)
	if err != nil {
		return err
	}

	// Look for history files
	historyPath := filepath.Join(historyDir, safeName)
	versions := []VersionInfo{}

	if entries, err := os.ReadDir(historyPath); err == nil {
		for _, entry := range entries {
			if strings.HasSuffix(entry.Name(), ".md") {
				// Parse version from filename: SKILL.v1.md, SKILL.v2.md, etc.
				parts := strings.Split(entry.Name(), ".v")
				if len(parts) >= 2 {
					var ver int
					fmt.Sscanf(parts[1], "%d.md", &ver)
					if ver > 0 {
						fi, _ := entry.Info()
						versions = append(versions, VersionInfo{
							Version:   ver,
							Timestamp: fi.ModTime(),
							Filename:  entry.Name(),
						})
					}
				}
			}
		}
	}

	// Sort by version descending
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].Version > versions[j].Version
	})

	fmt.Println()
	fmt.Printf("  Version History: %s\n", name)
	fmt.Println("  ════════════════════════════════════════════")
	fmt.Println()

	// Current version
	fmt.Printf("  ● current     %s  (active)\n", info.ModTime().Format("2006-01-02 15:04:05"))

	if len(versions) == 0 {
		fmt.Println()
		fmt.Println("  No previous versions found.")
		fmt.Println("  Versions are created automatically when you edit a skill.")
	} else {
		for _, v := range versions {
			fmt.Printf("  ○ v%-10d %s\n", v.Version, v.Timestamp.Format("2006-01-02 15:04:05"))
		}
	}

	fmt.Println()
	if len(versions) > 0 {
		fmt.Println("  To restore a version: openskill rollback", name, "<version>")
		fmt.Println()
	}

	return nil
}

// SaveVersion saves the current skill as a versioned backup
func SaveVersion(name string) error {
	safeName := strings.ReplaceAll(strings.ToLower(name), " ", "-")
	skillPath := filepath.Join(".claude/skills", safeName, "SKILL.md")

	// Read current content
	content, err := os.ReadFile(skillPath)
	if err != nil {
		return err
	}

	// Create history directory
	historyPath := filepath.Join(historyDir, safeName)
	if err := os.MkdirAll(historyPath, 0755); err != nil {
		return err
	}

	// Find next version number
	nextVersion := 1
	if entries, err := os.ReadDir(historyPath); err == nil {
		for _, entry := range entries {
			if strings.HasSuffix(entry.Name(), ".md") {
				parts := strings.Split(entry.Name(), ".v")
				if len(parts) >= 2 {
					var ver int
					fmt.Sscanf(parts[1], "%d.md", &ver)
					if ver >= nextVersion {
						nextVersion = ver + 1
					}
				}
			}
		}
	}

	// Save versioned copy
	versionFile := filepath.Join(historyPath, fmt.Sprintf("SKILL.v%d.md", nextVersion))
	return os.WriteFile(versionFile, content, 0644)
}
