package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"openskill/pkg/core"
	"openskill/pkg/skills"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var importFormat string
var importName string
var importImprove bool
var importAll bool

var ImportCmd = &cobra.Command{
	Use:   "import <source>",
	Short: "Import skills from GitHub repos, files, or URLs",
	Long: `Import skills from various sources:

- GitHub repositories (owner/repo or full URL)
- GitHub Gists
- Local files (.json, .yaml, .yml, .md)
- Raw URLs
- Stdin (use - as filename)

When importing from a GitHub repo, OpenSkill will look for SKILL.md files
in the repository and import them. Use --improve to enhance imported skills with AI.`,
	Args: cobra.ExactArgs(1),
	Example: `  openskill import anthropics/skills
  openskill import nummanali/openskills
  openskill import github.com/user/skills-repo
  openskill import skill.json
  openskill import https://gist.githubusercontent.com/.../skill.json
  openskill import anthropics/skills --improve  # Import and enhance with AI
  openskill import user/repo --all              # Import all skills from repo`,
	RunE: func(cmd *cobra.Command, args []string) error {
		source := args[0]

		// Check if it's a GitHub repo pattern (owner/repo)
		if isGitHubRepo(source) {
			return importFromGitHub(source)
		}

		// Original import logic for files/URLs
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
			} else if strings.HasSuffix(source, ".md") {
				detectedFormat = "md"
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
			case ".md":
				detectedFormat = "md"
			default:
				detectedFormat = importFormat
			}
		}

		if importFormat != "" {
			detectedFormat = importFormat
		}

		mgr := skills.NewManager()

		var skill *core.Skill
		var err error

		if detectedFormat == "md" {
			skill, err = parseSkillMD(content)
		} else {
			skill, err = mgr.Import(content, detectedFormat)
		}

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

		// Run improve if requested
		if importImprove {
			fmt.Println("\n  Running AI improvement...")
			improveCmd := exec.Command("openskill", "improve", skill.Name)
			improveCmd.Stdout = os.Stdout
			improveCmd.Stderr = os.Stderr
			if err := improveCmd.Run(); err != nil {
				fmt.Printf("  Warning: Could not improve skill: %v\n", err)
			}
		}

		return nil
	},
}

// isGitHubRepo checks if the source looks like a GitHub repo reference
func isGitHubRepo(source string) bool {
	// Match patterns like "owner/repo", "github.com/owner/repo"
	if strings.HasPrefix(source, "github.com/") {
		return true
	}
	// Simple owner/repo pattern (no dots, slashes only once)
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+/[a-zA-Z0-9_-]+$`, source)
	return matched
}

// importFromGitHub imports skills from a GitHub repository
func importFromGitHub(source string) error {
	// Normalize the source
	owner, repo := parseGitHubSource(source)
	if owner == "" || repo == "" {
		return fmt.Errorf("invalid GitHub repository format: %s", source)
	}

	fmt.Printf("Fetching skills from github.com/%s/%s...\n\n", owner, repo)

	// Get repository contents via GitHub API
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents", owner, repo)
	foundSkills, err := findSkillsInRepo(apiURL, "")
	if err != nil {
		return fmt.Errorf("failed to fetch repository: %w", err)
	}

	if len(foundSkills) == 0 {
		return fmt.Errorf("no SKILL.md files found in repository")
	}

	fmt.Printf("Found %d skill(s):\n", len(foundSkills))
	for _, s := range foundSkills {
		fmt.Printf("  - %s\n", s.path)
	}
	fmt.Println()

	// Import each skill
	mgr := skills.NewManager()
	importedCount := 0

	for _, si := range foundSkills {
		skill, err := fetchAndParseSkill(si.downloadURL)
		if err != nil {
			fmt.Printf("  ✗ Failed to import %s: %v\n", si.path, err)
			continue
		}

		if skill.Name == "" {
			skill.Name = si.name
		}

		// Check if skill already exists
		existing, _ := mgr.Get(skill.Name)
		if existing != nil && !importAll {
			fmt.Printf("  ⊘ Skipped %s (already exists, use --all to overwrite)\n", skill.Name)
			continue
		}

		if err := mgr.Add(skill); err != nil {
			fmt.Printf("  ✗ Failed to save %s: %v\n", skill.Name, err)
			continue
		}

		fmt.Printf("  ✓ Imported: %s\n", skill.Name)
		importedCount++

		// Run improve if requested
		if importImprove {
			fmt.Printf("    Improving with AI...\n")
			improveCmd := exec.Command("openskill", "improve", skill.Name)
			if err := improveCmd.Run(); err != nil {
				fmt.Printf("    Warning: Could not improve: %v\n", err)
			} else {
				fmt.Printf("    ✓ Enhanced with AI\n")
			}
		}
	}

	fmt.Printf("\nImported %d/%d skills\n", importedCount, len(foundSkills))
	return nil
}

func parseGitHubSource(source string) (owner, repo string) {
	// Remove github.com/ prefix if present
	source = strings.TrimPrefix(source, "github.com/")
	source = strings.TrimPrefix(source, "https://github.com/")

	parts := strings.Split(source, "/")
	if len(parts) >= 2 {
		return parts[0], parts[1]
	}
	return "", ""
}

type repoContent struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Type        string `json:"type"`
	DownloadURL string `json:"download_url"`
	URL         string `json:"url"`
}

type skillInfo struct {
	name        string
	path        string
	downloadURL string
}

func findSkillsInRepo(apiURL, basePath string) ([]skillInfo, error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	var contents []repoContent
	if err := json.NewDecoder(resp.Body).Decode(&contents); err != nil {
		return nil, err
	}

	var foundSkills []skillInfo

	for _, item := range contents {
		if item.Type == "file" && strings.ToUpper(item.Name) == "SKILL.MD" {
			// Extract skill name from directory
			skillName := filepath.Base(filepath.Dir(item.Path))
			if skillName == "." || skillName == "" {
				skillName = strings.TrimSuffix(item.Name, filepath.Ext(item.Name))
			}
			foundSkills = append(foundSkills, skillInfo{
				name:        skillName,
				path:        item.Path,
				downloadURL: item.DownloadURL,
			})
		} else if item.Type == "dir" {
			// Recursively search subdirectories
			subSkills, err := findSkillsInRepo(item.URL, item.Path)
			if err == nil {
				foundSkills = append(foundSkills, subSkills...)
			}
		}
	}

	return foundSkills, nil
}

func fetchAndParseSkill(url string) (*core.Skill, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return parseSkillMD(string(data))
}

// parseSkillMD parses a SKILL.md file into a Skill struct
func parseSkillMD(content string) (*core.Skill, error) {
	skill := &core.Skill{}

	// Check for YAML frontmatter
	if strings.HasPrefix(content, "---") {
		parts := strings.SplitN(content, "---", 3)
		if len(parts) >= 3 {
			// Parse YAML frontmatter
			if err := yaml.Unmarshal([]byte(parts[1]), skill); err != nil {
				return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
			}
			// The rest is the markdown content - extract rules from it
			mdContent := strings.TrimSpace(parts[2])
			skill.Rules = extractRulesFromMarkdown(mdContent)
		}
	} else {
		// No frontmatter, try to extract from markdown structure
		skill.Rules = extractRulesFromMarkdown(content)
		// Try to get name from first heading
		lines := strings.Split(content, "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "# ") {
				skill.Name = strings.TrimPrefix(line, "# ")
				break
			}
		}
	}

	return skill, nil
}

// extractRulesFromMarkdown extracts rules/instructions from markdown content
func extractRulesFromMarkdown(content string) []string {
	var rules []string
	lines := strings.Split(content, "\n")

	inRulesSection := false
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Check for rules/instructions section
		if strings.HasPrefix(strings.ToLower(line), "## rules") ||
			strings.HasPrefix(strings.ToLower(line), "## instructions") {
			inRulesSection = true
			continue
		}

		// End of section
		if inRulesSection && strings.HasPrefix(line, "## ") {
			break
		}

		// Extract list items as rules
		if inRulesSection {
			if strings.HasPrefix(line, "- ") {
				rules = append(rules, strings.TrimPrefix(line, "- "))
			} else if strings.HasPrefix(line, "* ") {
				rules = append(rules, strings.TrimPrefix(line, "* "))
			} else if matched, _ := regexp.MatchString(`^\d+\.\s+`, line); matched {
				// Numbered list
				re := regexp.MustCompile(`^\d+\.\s+`)
				rules = append(rules, re.ReplaceAllString(line, ""))
			}
		}
	}

	// If no rules section found, treat all list items as rules
	if len(rules) == 0 {
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "- ") {
				rules = append(rules, strings.TrimPrefix(line, "- "))
			} else if strings.HasPrefix(line, "* ") {
				rules = append(rules, strings.TrimPrefix(line, "* "))
			}
		}
	}

	return rules
}

func init() {
	ImportCmd.Flags().StringVarP(&importFormat, "format", "f", "yaml", "Import format (json, yaml, md)")
	ImportCmd.Flags().StringVarP(&importName, "name", "n", "", "Override skill name")
	ImportCmd.Flags().BoolVar(&importImprove, "improve", false, "Enhance imported skills with AI")
	ImportCmd.Flags().BoolVar(&importAll, "all", false, "Import all skills (overwrite existing)")
}
