# Contributing to OpenSkill

Thank you for your interest in contributing to OpenSkill! This document provides guidelines and steps for contributing.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Submitting a Pull Request](#submitting-a-pull-request)
- [Coding Standards](#coding-standards)
- [Commit Messages](#commit-messages)
- [Reporting Bugs](#reporting-bugs)
- [Suggesting Features](#suggesting-features)

## Code of Conduct

This project adheres to a [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/openskill.git
   cd openskill
   ```
3. **Add the upstream remote**:
   ```bash
   git remote add upstream https://github.com/rakshit-gen/openskill.git
   ```

## Development Setup

### Prerequisites

- Go 1.21 or later
- Make (optional, but recommended)
- A Groq API key for testing AI features

### Building

```bash
# Install dependencies
make deps

# Build the binary
make build

# Run tests
make test

# Run linter
make lint
```

### Running Locally

```bash
# After building, run from the project root
./openskill --help

# Or install to your PATH
sudo make install
openskill --help
```

## Making Changes

1. **Create a branch** for your changes:
   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b fix/your-bug-fix
   ```

2. **Make your changes** following the [coding standards](#coding-standards)

3. **Test your changes**:
   ```bash
   make test
   make lint
   ```

4. **Commit your changes** following the [commit message guidelines](#commit-messages)

## Submitting a Pull Request

1. **Push your branch** to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

2. **Create a Pull Request** on GitHub from your fork to the main repository

3. **Fill out the PR template** with:
   - A clear description of what your changes do
   - Any related issues (use "Fixes #123" to auto-close)
   - Screenshots if applicable (for UI changes)

4. **Wait for review** - maintainers will review your PR and may request changes

5. **Address feedback** - push additional commits to address review comments

## Coding Standards

### Go Style

- Follow the [Effective Go](https://go.dev/doc/effective_go) guidelines
- Use `gofmt` for formatting (run `make lint`)
- Keep functions small and focused
- Write descriptive variable and function names

### Project Structure

```
openskill/
├── cmd/openskill/          # CLI entry point and commands
│   ├── main.go
│   └── commands/
├── pkg/                    # Library code
│   ├── core/              # Core types and interfaces
│   ├── skills/            # Skill management
│   └── llm/               # LLM integration
└── test/                   # Integration tests
```

### Guidelines

- **Keep it simple** - Don't over-engineer solutions
- **Write tests** - New features should include tests
- **Document exports** - All exported functions need doc comments
- **Handle errors** - Always handle errors appropriately
- **No panics** - Return errors instead of panicking

## Commit Messages

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

### Types

- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

### Examples

```
feat(skills): add skill import from URL

fix(llm): handle rate limit errors gracefully

docs: update installation instructions for Linux ARM64

refactor(core): simplify skill validation logic
```

## Reporting Bugs

Before reporting a bug:

1. **Search existing issues** to avoid duplicates
2. **Try the latest version** - your bug might be fixed

When reporting:

1. Use the **Bug Report** issue template
2. Include:
   - OpenSkill version (`openskill --version`)
   - Operating system and architecture
   - Steps to reproduce
   - Expected vs actual behavior
   - Error messages or logs

## Suggesting Features

We welcome feature suggestions! Before suggesting:

1. **Check the roadmap** in the README
2. **Search existing issues** for similar suggestions

When suggesting:

1. Use the **Feature Request** issue template
2. Describe:
   - The problem you're trying to solve
   - Your proposed solution
   - Any alternatives you've considered

## Questions?

- Open a [Discussion](https://github.com/rakshit-gen/openskill/discussions) for general questions
- Join the community chat (coming soon)

---

Thank you for contributing to OpenSkill!
