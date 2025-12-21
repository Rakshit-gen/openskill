package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"openskill/pkg/config"

	"github.com/spf13/cobra"
)

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage OpenSkill configuration",
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> [value]",
	Short: "Set a configuration value",
	Long: `Set a configuration value. Supported keys:
  api-key    Your Groq API key
  model      LLM model to use (default: llama-3.3-70b-versatile)

If value is not provided, you will be prompted to enter it (useful for secrets).`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		var value string

		if len(args) == 2 {
			value = args[1]
		} else {
			// Prompt for value (useful for API keys)
			fmt.Printf("Enter value for %s: ", key)
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			value = strings.TrimSpace(input)
		}

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		switch key {
		case "api-key":
			cfg.GroqAPIKey = value
		case "model":
			cfg.Model = value
		default:
			return fmt.Errorf("unknown config key: %s (supported: api-key, model)", key)
		}

		if err := config.Save(cfg); err != nil {
			return err
		}

		fmt.Printf("✓ Set %s\n", key)
		return nil
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		switch key {
		case "api-key":
			if cfg.GroqAPIKey == "" {
				fmt.Println("(not set)")
			} else {
				// Mask the API key for security
				masked := cfg.GroqAPIKey[:4] + "..." + cfg.GroqAPIKey[len(cfg.GroqAPIKey)-4:]
				fmt.Println(masked)
			}
		case "model":
			if cfg.Model == "" {
				fmt.Println("llama-3.3-70b-versatile (default)")
			} else {
				fmt.Println(cfg.Model)
			}
		default:
			return fmt.Errorf("unknown config key: %s (supported: api-key, model)", key)
		}

		return nil
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration values",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		fmt.Println("OpenSkill Configuration:")
		fmt.Println()

		// API Key
		if cfg.GroqAPIKey == "" {
			fmt.Println("  api-key: (not set)")
		} else {
			masked := cfg.GroqAPIKey[:4] + "..." + cfg.GroqAPIKey[len(cfg.GroqAPIKey)-4:]
			fmt.Printf("  api-key: %s\n", masked)
		}

		// Model
		if cfg.Model == "" {
			fmt.Println("  model:   llama-3.3-70b-versatile (default)")
		} else {
			fmt.Printf("  model:   %s\n", cfg.Model)
		}

		fmt.Println()
		fmt.Println("Config file: ~/.openskill/config.yaml")
		fmt.Println("Note: Environment variables (GROQ_API_KEY, OPENSKILL_MODEL) take precedence")

		return nil
	},
}

func init() {
	ConfigCmd.AddCommand(configSetCmd)
	ConfigCmd.AddCommand(configGetCmd)
	ConfigCmd.AddCommand(configListCmd)
}
