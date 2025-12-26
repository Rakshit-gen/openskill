<h1 align="center">OpenSkill CLI</h1>

<p align="center">
  <strong>Create and manage Claude skills with AI-powered content generation</strong>
</p>

<p align="center">
  <a href="https://github.com/rakshit-gen/openskill/releases"><img src="https://img.shields.io/github/v/release/rakshit-gen/openskill?style=flat-square" alt="Release"></a>
  <a href="https://github.com/rakshit-gen/openskill/blob/main/LICENSE"><img src="https://img.shields.io/github/license/rakshit-gen/openskill?style=flat-square" alt="License"></a>
  <a href="https://github.com/rakshit-gen/openskill/stargazers"><img src="https://img.shields.io/github/stars/rakshit-gen/openskill?style=flat-square" alt="Stars"></a>
  <a href="https://github.com/rakshit-gen/openskill/issues"><img src="https://img.shields.io/github/issues/rakshit-gen/openskill?style=flat-square" alt="Issues"></a>
</p>

<p align="center">
  <a href="#installation">Installation</a> •
  <a href="#quick-start">Quick Start</a> •
  <a href="#commands">Commands</a> •
  <a href="#configuration">Configuration</a> •
  <a href="#contributing">Contributing</a>
</p>

---

## What is OpenSkill?

OpenSkill is a command-line tool that simplifies creating and managing [Claude](https://claude.ai) skills. It uses AI to automatically generate comprehensive skill descriptions and rules from simple prompts.

### Features

- **Multi-Provider AI** - Choose from Groq, OpenAI, Anthropic, or Ollama (local)
- **SKILL.md Format** - Skills stored as Markdown with YAML frontmatter for Claude's native skill discovery
- **Version History** - Track changes with automatic versioning and rollback support
- **Skill Composition** - Extend and combine skills with `extends` and `includes`
- **Validation** - Validate skill structure before deploying
- **Fast** - Built in Go for maximum performance

### Supported AI Providers

| Provider | Model (Default) | API Key |
|----------|-----------------|---------|
| **Groq** | llama-3.3-70b-versatile | [console.groq.com](https://console.groq.com) |
| **OpenAI** | gpt-4o-mini | [platform.openai.com](https://platform.openai.com) |
| **Anthropic** | claude-3-5-sonnet-20241022 | [console.anthropic.com](https://console.anthropic.com) |
| **Ollama** | llama3.2 | No API key (runs locally) |

## Installation

### Quick Install (Recommended)

```bash
curl -fsSL openskill.online/api/install | bash
```

### Build from Source

Requires Go 1.21+

```bash
git clone https://github.com/rakshit-gen/openskill.git
cd openskill && make build && sudo make install
```

### Verify Installation

```bash
openskill --help
```

## Configuration

### Set Up Your AI Provider

OpenSkill supports multiple AI providers. Set your preferred provider and API key:

```bash
# Set provider (default: groq)
openskill config set provider groq    # or: openai, anthropic, ollama

# Set API key for your provider
openskill config set api-key
```

### Provider-Specific Setup

#### Groq (Default - Free & Fast)

```bash
openskill config set provider groq
openskill config set groq-api-key YOUR_KEY
```

#### OpenAI

```bash
openskill config set provider openai
openskill config set openai-api-key YOUR_KEY
```

#### Anthropic

```bash
openskill config set provider anthropic
openskill config set anthropic-api-key YOUR_KEY
```

#### Ollama (Local - No API Key)

```bash
# Make sure Ollama is running: ollama serve
openskill config set provider ollama
openskill config set ollama-model llama3.2
```

### View Configuration

```bash
openskill config list
```

### Environment Variables

Environment variables take precedence over config file:

| Variable | Description |
|----------|-------------|
| `OPENSKILL_PROVIDER` | Active provider (groq, openai, anthropic, ollama) |
| `GROQ_API_KEY` | Groq API key |
| `OPENAI_API_KEY` | OpenAI API key |
| `ANTHROPIC_API_KEY` | Anthropic API key |
| `OPENSKILL_MODEL` | Override model for any provider |
| `OLLAMA_HOST` | Custom Ollama endpoint |

## Quick Start

### Initialize OpenSkill

```bash
openskill init
```

### Create your first skill

```bash
openskill add "code-review" -d "Reviews code for best practices"
```

The AI will generate a comprehensive skill:

```
Generating skill with Groq...

✓ Added skill: code-review
  Description: Comprehensive code review focusing on security, performance, and maintainability
  Rules:
    1. Check for security vulnerabilities (XSS, SQL injection, etc.)
    2. Verify proper error handling and edge cases
    3. Ensure code follows project conventions
    4. Review test coverage and quality
```

### List all skills

```bash
openskill list
```

### View skill details

```bash
openskill show "code-review"
```

## Commands

| Command | Description |
|---------|-------------|
| `openskill init` | Initialize OpenSkill in your project |
| `openskill add <name> -d <description>` | Create a new skill with AI generation |
| `openskill add <name> -d <desc> --manual -r <rule>` | Create skill manually with custom rules |
| `openskill list` | List all skills |
| `openskill show <name>` | Show detailed skill information |
| `openskill edit <name> -d <description>` | Update skill description |
| `openskill edit <name> -r <rule1> -r <rule2>` | Replace skill rules |
| `openskill remove <name>` | Delete a skill |
| `openskill validate <name>` | Validate skill structure |
| `openskill history <name>` | Show version history |
| `openskill rollback <name> <version>` | Restore a previous version |
| `openskill config set <key> [value]` | Set configuration |
| `openskill config get <key>` | Get configuration value |
| `openskill config list` | List all configuration |

### Flags

| Flag | Description |
|------|-------------|
| `-d, --desc` | Skill description (required for add) |
| `-r, --rule` | Add a rule (can be repeated, manual mode only) |
| `--manual` | Skip AI generation, use provided values |

## Skill Format

Skills are stored as Markdown files with YAML frontmatter in `.claude/skills/<name>/SKILL.md`:

```
.claude/
└── skills/
    ├── code-review/
    │   └── SKILL.md
    └── bug-finder/
        └── SKILL.md
```

### SKILL.md Structure

```markdown
---
name: code-review
description: Comprehensive code review focusing on security and maintainability
---

# code-review

Comprehensive code review focusing on security, performance,
and maintainability best practices.

## Rules

- Check for security vulnerabilities (XSS, SQL injection, etc.)
- Verify proper error handling and edge cases
- Ensure code follows project conventions
- Review test coverage and quality
```

### Skill Composition

Extend skills with `extends`:

```markdown
---
name: security-review
description: Security-focused code review
extends: code-review
---

# security-review

## Rules

- Focus on OWASP Top 10 vulnerabilities
- Check authentication and authorization flows
```

Combine skills with `includes`:

```markdown
---
name: full-review
description: Comprehensive review combining multiple aspects
includes:
  - code-review
  - security-review
  - performance-review
---

# full-review

## Rules

- Provide a summary score for each review area
```

## Project Structure

```
openskill/
├── cmd/
│   └── openskill/
│       ├── main.go           # Entry point
│       └── commands/         # CLI commands
│           ├── init.go       # Initialize project
│           ├── add.go        # Add skills
│           ├── list.go       # List skills
│           ├── show.go       # Show skill details
│           ├── edit.go       # Edit skills
│           ├── remove.go     # Remove skills
│           ├── validate.go   # Validate skills
│           ├── history.go    # Version history
│           ├── rollback.go   # Rollback versions
│           └── config.go     # Configuration
├── pkg/
│   ├── core/
│   │   └── skill.go          # Skill data structure
│   ├── skills/
│   │   └── manager.go        # Skill file management
│   ├── llm/
│   │   ├── provider.go       # Provider interface
│   │   ├── generator.go      # AI generation
│   │   ├── groq.go           # Groq client
│   │   ├── openai.go         # OpenAI client
│   │   ├── anthropic.go      # Anthropic client
│   │   └── ollama.go         # Ollama client
│   └── config/
│       └── config.go         # Configuration management
├── Makefile
├── go.mod
└── go.sum
```

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Setup

1. Fork and clone the repository
2. Install Go 1.21+
3. Run `make deps` to install dependencies
4. Run `make build` to build
5. Run `make test` to run tests

### Quick Commands

```bash
make build      # Build binary
make install    # Install to /usr/local/bin
make test       # Run tests
make clean      # Clean build artifacts
make deps       # Tidy dependencies
make lint       # Run linter
make release    # Build for all platforms
```

## License

MIT License - see [LICENSE](LICENSE) for details.

## Acknowledgments

- [Groq](https://groq.com) for fast LLM inference
- [OpenAI](https://openai.com) for GPT models
- [Anthropic](https://anthropic.com) for Claude models
- [Ollama](https://ollama.ai) for local LLM support
- [Cobra](https://github.com/spf13/cobra) for CLI framework
- [Claude](https://claude.ai) for inspiration

---

<p align="center">
  Made with ❤️ by <a href="https://github.com/rakshit-gen">Rakshit-gen</a>
</p>
