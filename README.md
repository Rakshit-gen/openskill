<p align="center">
  <img src="assets/logo.svg" width="120" alt="OpenSkill Logo">
</p>

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

OpenSkill is a command-line tool that simplifies creating and managing [Claude](https://claude.ai) skills. It uses AI (powered by [Groq](https://groq.com)) to automatically generate comprehensive skill descriptions and rules from simple prompts.

### Features

- **AI-Powered Generation** - Automatically generate detailed skill content from brief descriptions
- **Simple CLI** - Intuitive commands that get out of your way
- **YAML Storage** - Skills stored as readable YAML files for easy version control
- **Local First** - Your skills stay on your machine, no cloud sync required
- **Fast** - Built in Go for maximum performance

## Installation

### Quick Install (Recommended)

#### macOS (Apple Silicon)

```bash
curl -L https://github.com/rakshit-gen/openskill/releases/download/v0.1.0/openskill_darwin_arm64.tar.gz | tar xz
sudo mv openskill /usr/local/bin/
```

#### macOS (Intel)

```bash
curl -L https://github.com/rakshit-gen/openskill/releases/download/v0.1.0/openskill_darwin_amd64.tar.gz | tar xz
sudo mv openskill /usr/local/bin/
```

#### Linux (x86_64)

```bash
curl -L https://github.com/rakshit-gen/openskill/releases/download/v0.1.0/openskill_linux_amd64.tar.gz | tar xz
sudo mv openskill /usr/local/bin/
```

#### Linux (ARM64)

```bash
curl -L https://github.com/rakshit-gen/openskill/releases/download/v0.1.0/openskill_linux_arm64.tar.gz | tar xz
sudo mv openskill /usr/local/bin/
```

### Build from Source

Requires Go 1.21+

```bash
git clone https://github.com/rakshit-gen/openskill.git
cd openskill
make build
sudo make install
```

### Verify Installation

```bash
openskill --help
```

## Configuration

### Groq API Key

OpenSkill uses Groq's LLM for AI-powered skill generation. Get your free API key:

1. Visit [console.groq.com](https://console.groq.com)
2. Create a free account
3. Generate an API key

Set your API key:

```bash
export GROQ_API_KEY=your_key_here
```

Add this to your shell profile (`~/.bashrc`, `~/.zshrc`, etc.) for persistence.

## Quick Start

### Create your first skill

```bash
openskill add "code-review" -d "Reviews code for best practices"
```

The AI will generate a comprehensive skill with detailed rules:

```
Generating skill with AI...

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
| `openskill add <name> -d <description>` | Create a new skill with AI generation |
| `openskill add <name> -d <desc> --manual -r <rule>` | Create skill manually with custom rules |
| `openskill list` | List all skills |
| `openskill show <name>` | Show detailed skill information |
| `openskill edit <name> -d <description>` | Update skill description |
| `openskill edit <name> -r <rule1> -r <rule2>` | Replace skill rules |
| `openskill remove <name>` | Delete a skill |

### Flags

| Flag | Description |
|------|-------------|
| `-d, --desc` | Skill description (required for add) |
| `-r, --rule` | Add a rule (can be repeated, manual mode only) |
| `--manual` | Skip AI generation, use provided values |

## Skill Format

Skills are stored as YAML files in `.claude/skills/` directory:

```yaml
name: code-review
description: >
  Comprehensive code review focusing on security,
  performance, and maintainability best practices.
rules:
  - Check for security vulnerabilities
  - Verify proper error handling
  - Ensure code follows conventions
  - Review test coverage
```

## Project Structure

```
openskill/
├── cmd/
│   └── openskill/
│       ├── main.go           # Entry point
│       └── commands/         # CLI commands
│           ├── add.go
│           ├── list.go
│           ├── show.go
│           ├── edit.go
│           └── remove.go
├── pkg/
│   ├── core/
│   │   └── skill.go          # Skill data structure
│   ├── skills/
│   │   └── manager.go        # Skill file management
│   └── llm/
│       ├── generator.go      # AI generation interface
│       └── groq.go           # Groq API client
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

## Roadmap

- [ ] Support for multiple LLM providers (OpenAI, Anthropic, Ollama)
- [ ] Skill templates and sharing
- [ ] Interactive skill editor
- [ ] Skill import/export
- [ ] Plugin system

## License

MIT License - see [LICENSE](LICENSE) for details.

## Acknowledgments

- [Groq](https://groq.com) for fast LLM inference
- [Cobra](https://github.com/spf13/cobra) for CLI framework
- [Claude](https://claude.ai) for inspiration

---

<p align="center">
  Made with ❤️ by <a href="https://github.com/rakshit-gen">Rakshit-gen</a>
</p>
