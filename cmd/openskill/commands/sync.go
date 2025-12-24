package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"openskill/pkg/skills"

	"github.com/spf13/cobra"
)

var syncRemote string
var syncPush bool
var syncPull bool

var SyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync skills with a remote repository",
	Long: `Synchronize skills with a Git repository for backup and sharing.

This command helps you:
- Back up skills to a remote repository
- Share skills across machines
- Collaborate on skills with team members`,
	Example: `  openskill sync --remote git@github.com:user/skills.git
  openskill sync --push
  openskill sync --pull`,
	RunE: func(cmd *cobra.Command, args []string) error {
		skillsDir := skills.SkillsDir

		// Check if .claude/skills is a git repo
		gitDir := filepath.Join(skillsDir, ".git")
		isGitRepo := false
		if _, err := os.Stat(gitDir); err == nil {
			isGitRepo = true
		}

		if syncRemote != "" {
			// Initialize or update remote
			if !isGitRepo {
				fmt.Println("Initializing git repository in .claude/skills...")
				if err := runGitCommand(skillsDir, "init"); err != nil {
					return fmt.Errorf("git init failed: %w", err)
				}
			}

			// Add or update remote
			if err := runGitCommand(skillsDir, "remote", "remove", "origin"); err != nil {
				// Ignore error if remote doesn't exist
			}
			if err := runGitCommand(skillsDir, "remote", "add", "origin", syncRemote); err != nil {
				return fmt.Errorf("failed to add remote: %w", err)
			}

			fmt.Printf("✓ Remote set to: %s\n", syncRemote)
			return nil
		}

		if !isGitRepo {
			return fmt.Errorf("skills directory is not a git repository. Use --remote to initialize")
		}

		if syncPush {
			// Add all changes
			if err := runGitCommand(skillsDir, "add", "-A"); err != nil {
				return fmt.Errorf("git add failed: %w", err)
			}

			// Check if there are changes to commit
			status, _ := getGitOutput(skillsDir, "status", "--porcelain")
			if strings.TrimSpace(status) != "" {
				// Commit changes
				if err := runGitCommand(skillsDir, "commit", "-m", "Update skills"); err != nil {
					return fmt.Errorf("git commit failed: %w", err)
				}
			}

			// Push to remote
			fmt.Println("Pushing skills to remote...")
			if err := runGitCommand(skillsDir, "push", "-u", "origin", "main"); err != nil {
				// Try master branch
				if err := runGitCommand(skillsDir, "push", "-u", "origin", "master"); err != nil {
					return fmt.Errorf("git push failed: %w", err)
				}
			}

			fmt.Println("✓ Skills pushed to remote")
			return nil
		}

		if syncPull {
			fmt.Println("Pulling skills from remote...")
			if err := runGitCommand(skillsDir, "pull", "origin", "main"); err != nil {
				// Try master branch
				if err := runGitCommand(skillsDir, "pull", "origin", "master"); err != nil {
					return fmt.Errorf("git pull failed: %w", err)
				}
			}

			fmt.Println("✓ Skills pulled from remote")
			return nil
		}

		// Default: show status
		fmt.Println("Sync Status:")
		fmt.Println("───────────────────────────────────")

		// Get remote URL
		remote, _ := getGitOutput(skillsDir, "remote", "get-url", "origin")
		if remote != "" {
			fmt.Printf("Remote: %s\n", strings.TrimSpace(remote))
		} else {
			fmt.Println("Remote: Not configured")
		}

		// Get current status
		status, _ := getGitOutput(skillsDir, "status", "--short")
		if strings.TrimSpace(status) != "" {
			fmt.Println("\nUnsynced changes:")
			fmt.Println(status)
		} else {
			fmt.Println("\nAll skills are synced.")
		}

		fmt.Println("\nUse --push to upload changes, --pull to download.")

		return nil
	},
}

func runGitCommand(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func getGitOutput(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	output, err := cmd.Output()
	return string(output), err
}

func init() {
	SyncCmd.Flags().StringVar(&syncRemote, "remote", "", "Set the remote repository URL")
	SyncCmd.Flags().BoolVar(&syncPush, "push", false, "Push local changes to remote")
	SyncCmd.Flags().BoolVar(&syncPull, "pull", false, "Pull remote changes to local")
}
