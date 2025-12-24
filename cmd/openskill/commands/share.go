package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"openskill/pkg/core"
	"openskill/pkg/skills"

	"github.com/spf13/cobra"
)

var sharePublic bool
var shareDescription string

// ShareResponse represents the API response from the marketplace
type ShareResponse struct {
	Success bool   `json:"success"`
	URL     string `json:"url"`
	ID      string `json:"id"`
	Message string `json:"message"`
}

var ShareCmd = &cobra.Command{
	Use:   "share <skill-name>",
	Short: "Share a skill to the OpenSkill marketplace",
	Long: `Upload a skill to openskill.dev marketplace for others to discover and use.

By default, skills are shared publicly. Use --public=false for unlisted sharing
(accessible via direct link only).

After sharing, you'll receive:
- A shareable URL (openskill.dev/skills/username/skill-name)
- An install command others can use`,
	Args: cobra.ExactArgs(1),
	Example: `  openskill share code-review
  openskill share my-skill --description "Custom code review for React"
  openskill share private-skill --public=false`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		mgr := skills.NewManager()

		// Get the skill
		skill, err := mgr.Get(name)
		if err != nil {
			return fmt.Errorf("skill not found: %s", name)
		}

		// Override description if provided
		if shareDescription != "" {
			skill.Description = shareDescription
		}

		// Get author from git config or environment
		author := getAuthor()
		if skill.Author == "" {
			skill.Author = author
		}

		fmt.Printf("Sharing skill: %s\n", skill.Name)
		fmt.Printf("  Author: %s\n", skill.Author)
		fmt.Printf("  Description: %s\n", skill.Description)
		fmt.Printf("  Rules: %d\n", len(skill.Rules))
		if sharePublic {
			fmt.Printf("  Visibility: Public\n")
		} else {
			fmt.Printf("  Visibility: Unlisted\n")
		}
		fmt.Println()

		// Prepare the payload
		payload := map[string]interface{}{
			"skill":  skill,
			"public": sharePublic,
			"author": skill.Author,
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to serialize skill: %w", err)
		}

		// For now, simulate the upload since we don't have a backend yet
		// In production, this would POST to the actual API
		apiURL := "https://api.openskill.dev/v1/skills"

		// Check if we're in development mode
		if os.Getenv("OPENSKILL_API_URL") != "" {
			apiURL = os.Getenv("OPENSKILL_API_URL") + "/v1/skills"
		}

		// Try to upload
		resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			// API not available - provide local export instead
			fmt.Println("  Note: Marketplace API not yet available.")
			fmt.Println("  Generating shareable export instead...\n")

			return generateLocalShare(skill, mgr)
		}
		defer resp.Body.Close()

		var shareResp ShareResponse
		if err := json.NewDecoder(resp.Body).Decode(&shareResp); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}

		if !shareResp.Success {
			return fmt.Errorf("upload failed: %s", shareResp.Message)
		}

		fmt.Printf("✓ Skill shared successfully!\n\n")
		fmt.Printf("  URL: %s\n", shareResp.URL)
		fmt.Printf("  Install: openskill import %s\n", shareResp.ID)

		return nil
	},
}

func getAuthor() string {
	// Try git config first
	if author := getGitConfig("user.name"); author != "" {
		return author
	}
	// Fall back to environment
	if author := os.Getenv("USER"); author != "" {
		return author
	}
	return "anonymous"
}

func getGitConfig(key string) string {
	// Simple git config lookup
	cmd := exec.Command("git", "config", "--get", key)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func generateLocalShare(skill *core.Skill, mgr *skills.Manager) error {
	// Export to YAML
	content, err := mgr.Export(skill.Name, "yaml")
	if err != nil {
		return err
	}

	// Generate a shareable filename
	filename := fmt.Sprintf("%s.skill.yaml", skill.Name)
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Printf("✓ Exported to: %s\n\n", filename)
	fmt.Println("Share options:")
	fmt.Println("  1. Upload to GitHub Gist:")
	fmt.Printf("     gh gist create %s --public\n\n", filename)
	fmt.Println("  2. Share the file directly")
	fmt.Println("     Others can import with: openskill import <url-or-file>")
	fmt.Println()
	fmt.Println("  3. Create a GitHub repo for your skills:")
	fmt.Println("     mkdir my-skills && mv *.skill.yaml my-skills/")
	fmt.Println("     cd my-skills && git init && gh repo create --public")
	fmt.Println("     Others can import with: openskill import username/my-skills")

	return nil
}

func init() {
	ShareCmd.Flags().BoolVar(&sharePublic, "public", true, "Make skill publicly discoverable")
	ShareCmd.Flags().StringVarP(&shareDescription, "description", "d", "", "Override skill description")
}
