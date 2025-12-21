package llm

// Provider represents an LLM provider interface
type Provider interface {
	Name() string
	Generate(prompt string) (string, error)
	IsConfigured() bool
}

// ProviderType represents the type of LLM provider
type ProviderType string

const (
	ProviderGroq     ProviderType = "groq"
	ProviderOpenAI   ProviderType = "openai"
	ProviderAnthropic ProviderType = "anthropic"
	ProviderOllama   ProviderType = "ollama"
)

// DefaultModels for each provider
var DefaultModels = map[ProviderType]string{
	ProviderGroq:     "llama-3.3-70b-versatile",
	ProviderOpenAI:   "gpt-4o-mini",
	ProviderAnthropic: "claude-3-5-sonnet-20241022",
	ProviderOllama:   "llama3.2",
}

// ProviderEndpoints for each provider
var ProviderEndpoints = map[ProviderType]string{
	ProviderGroq:     "https://api.groq.com/openai/v1/chat/completions",
	ProviderOpenAI:   "https://api.openai.com/v1/chat/completions",
	ProviderAnthropic: "https://api.anthropic.com/v1/messages",
	ProviderOllama:   "http://localhost:11434/api/chat",
}
