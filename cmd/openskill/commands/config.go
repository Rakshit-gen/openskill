package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"openskill/pkg/config"
	"openskill/pkg/llm"

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

Provider Selection:
  provider           Active AI provider (groq, openai, anthropic, ollama)

API Keys:
  api-key            API key for current provider
  groq-api-key       Groq API key
  openai-api-key     OpenAI API key
  anthropic-api-key  Anthropic API key

Models:
  model              Default model for all providers
  groq-model         Groq-specific model (default: llama-3.3-70b-versatile)
  openai-model       OpenAI-specific model (default: gpt-4o-mini)
  anthropic-model    Anthropic-specific model (default: claude-3-5-sonnet-20241022)
  ollama-model       Ollama-specific model (default: llama3.2)

Ollama:
  ollama-endpoint    Custom Ollama endpoint (default: http://localhost:11434)

If value is not provided, you will be prompted to enter it (useful for secrets).`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		var value string

		if len(args) == 2 {
			value = args[1]
		} else {
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
		case "provider":
			validProviders := []string{"groq", "openai", "anthropic", "ollama"}
			value = strings.ToLower(value)
			valid := false
			for _, p := range validProviders {
				if value == p {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("invalid provider: %s (valid: %s)", value, strings.Join(validProviders, ", "))
			}
			cfg.Provider = value

		case "api-key":
			provider := config.GetProvider()
			switch provider {
			case "groq":
				cfg.GroqAPIKey = value
			case "openai":
				cfg.OpenAIAPIKey = value
			case "anthropic":
				cfg.AnthropicAPIKey = value
			default:
				cfg.GroqAPIKey = value
			}
		case "groq-api-key":
			cfg.GroqAPIKey = value
		case "openai-api-key":
			cfg.OpenAIAPIKey = value
		case "anthropic-api-key":
			cfg.AnthropicAPIKey = value

		case "model":
			cfg.Model = value
		case "groq-model":
			cfg.GroqModel = value
		case "openai-model":
			cfg.OpenAIModel = value
		case "anthropic-model":
			cfg.AnthropicModel = value
		case "ollama-model":
			cfg.OllamaModel = value

		case "ollama-endpoint":
			cfg.OllamaEndpoint = value

		default:
			return fmt.Errorf("unknown config key: %s\nRun 'openskill config set --help' for available keys", key)
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
		case "provider":
			fmt.Println(config.GetProvider())

		case "api-key":
			apiKey := config.GetAPIKey()
			if apiKey == "" {
				fmt.Println("(not set)")
			} else {
				fmt.Println(maskKey(apiKey))
			}
		case "groq-api-key":
			if cfg.GroqAPIKey == "" {
				fmt.Println("(not set)")
			} else {
				fmt.Println(maskKey(cfg.GroqAPIKey))
			}
		case "openai-api-key":
			if cfg.OpenAIAPIKey == "" {
				fmt.Println("(not set)")
			} else {
				fmt.Println(maskKey(cfg.OpenAIAPIKey))
			}
		case "anthropic-api-key":
			if cfg.AnthropicAPIKey == "" {
				fmt.Println("(not set)")
			} else {
				fmt.Println(maskKey(cfg.AnthropicAPIKey))
			}

		case "model":
			fmt.Println(config.GetModel())
		case "groq-model":
			fmt.Println(config.GetProviderModel("groq"))
		case "openai-model":
			fmt.Println(config.GetProviderModel("openai"))
		case "anthropic-model":
			fmt.Println(config.GetProviderModel("anthropic"))
		case "ollama-model":
			fmt.Println(config.GetProviderModel("ollama"))

		case "ollama-endpoint":
			fmt.Println(config.GetOllamaEndpoint())

		default:
			return fmt.Errorf("unknown config key: %s", key)
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

		fmt.Println()
		fmt.Println("  OpenSkill Configuration")
		fmt.Println("  ═══════════════════════════════════════════════════════")
		fmt.Println()

		provider := config.GetProvider()
		fmt.Printf("  Provider:          %s\n", provider)
		fmt.Println()

		fmt.Println("  API Keys:")
		printConfigValue("    Groq", cfg.GroqAPIKey, true)
		printConfigValue("    OpenAI", cfg.OpenAIAPIKey, true)
		printConfigValue("    Anthropic", cfg.AnthropicAPIKey, true)
		fmt.Println()

		fmt.Println("  Models:")
		fmt.Printf("    Groq:            %s\n", config.GetProviderModel("groq"))
		fmt.Printf("    OpenAI:          %s\n", config.GetProviderModel("openai"))
		fmt.Printf("    Anthropic:       %s\n", config.GetProviderModel("anthropic"))
		fmt.Printf("    Ollama:          %s\n", config.GetProviderModel("ollama"))
		fmt.Println()

		fmt.Println("  Ollama:")
		fmt.Printf("    Endpoint:        %s\n", config.GetOllamaEndpoint())
		fmt.Println()

		available := llm.GetAvailableProviders()
		fmt.Printf("  Configured:        %s\n", strings.Join(available, ", "))
		fmt.Println()

		fmt.Println("  Config file:       ~/.openskill/config.yaml")
		fmt.Println()
		fmt.Println("  Environment variables (take precedence):")
		fmt.Println("    OPENSKILL_PROVIDER, GROQ_API_KEY, OPENAI_API_KEY,")
		fmt.Println("    ANTHROPIC_API_KEY, OPENSKILL_MODEL, OLLAMA_HOST")
		fmt.Println()

		return nil
	},
}

func maskKey(key string) string {
	if len(key) <= 8 {
		return "***"
	}
	return key[:4] + "..." + key[len(key)-4:]
}

func printConfigValue(label, value string, isSecret bool) {
	if value == "" {
		fmt.Printf("%s:       (not set)\n", label)
	} else if isSecret {
		fmt.Printf("%s:       %s\n", label, maskKey(value))
	} else {
		fmt.Printf("%s:       %s\n", label, value)
	}
}

func init() {
	ConfigCmd.AddCommand(configSetCmd)
	ConfigCmd.AddCommand(configGetCmd)
	ConfigCmd.AddCommand(configListCmd)
}
