# OpenSkill CLI

A CLI tool to create and manage Claude skills with AI-powered content generation.

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap rakshitsisodiya/tap
brew install openskill
```

### From Source

```bash
git clone https://github.com/rakshitsisodiya/openskill.git
cd openskill
make build
./build/openskill --help
```

## Setup

Get a free API key from [Groq](https://console.groq.com/) and set it:

```bash
export GROQ_API_KEY=your_key_here
```

## Usage

### Add a skill (AI-powered)

```bash
openskill add "code-review" -d "Reviews code for best practices"
```

The AI will generate a detailed description and relevant rules automatically.

### Add a skill (manual)

```bash
openskill add "code-review" -d "Reviews code" --manual -r "Check security" -r "Check performance"
```

### List skills

```bash
openskill list
```

### Show skill details

```bash
openskill show "code-review"
```

### Edit a skill

```bash
openskill edit "code-review" -d "New description"
openskill edit "code-review" -r "New rule 1" -r "New rule 2"
```

### Remove a skill

```bash
openskill remove "code-review"
```

## Where are skills stored?

Skills are saved to `.claude/skills/` in your current directory as YAML files.

## License

MIT
# openskill
