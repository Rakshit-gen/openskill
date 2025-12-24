package skills

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"openskill/pkg/core"

	"gopkg.in/yaml.v3"
)

const SkillsDir = ".claude/skills"
const TemplatesDir = ".claude/templates"
const GroupsDir = ".claude/groups"
const WorkspaceFile = ".claude/workspace.yaml"
const HistoryDir = ".claude/skills/.history"

// Manager handles skill CRUD operations
type Manager struct {
	baseDir string
}

// NewManager creates a new skill manager
func NewManager() *Manager {
	return &Manager{baseDir: SkillsDir}
}

// ensureDir creates the skills directory if it doesn't exist
func (m *Manager) ensureDir() error {
	return os.MkdirAll(m.baseDir, 0755)
}

// skillDir returns the directory path for a skill
func (m *Manager) skillDir(name string) string {
	safeName := strings.ReplaceAll(strings.ToLower(name), " ", "-")
	return filepath.Join(m.baseDir, safeName)
}

// skillPath returns the SKILL.md file path for a skill
func (m *Manager) skillPath(name string) string {
	return filepath.Join(m.skillDir(name), "SKILL.md")
}

// Add creates a new skill
func (m *Manager) Add(skill *core.Skill) error {
	if err := m.ensureDir(); err != nil {
		return fmt.Errorf("failed to create skills directory: %w", err)
	}

	dir := m.skillDir(skill.Name)
	if _, err := os.Stat(dir); err == nil {
		return fmt.Errorf("skill '%s' already exists", skill.Name)
	}

	// Create skill directory
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create skill directory: %w", err)
	}

	return m.save(skill)
}

// List returns all skills
func (m *Manager) List() ([]core.Skill, error) {
	if _, err := os.Stat(m.baseDir); os.IsNotExist(err) {
		return []core.Skill{}, nil
	}

	entries, err := os.ReadDir(m.baseDir)
	if err != nil {
		return nil, err
	}

	var skills []core.Skill
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Skip history directory
		if entry.Name() == ".history" {
			continue
		}

		// Check if SKILL.md exists in the directory
		skillPath := filepath.Join(m.baseDir, entry.Name(), "SKILL.md")
		if _, err := os.Stat(skillPath); os.IsNotExist(err) {
			continue
		}

		skill, err := m.load(entry.Name())
		if err != nil {
			continue
		}
		skills = append(skills, *skill)
	}

	return skills, nil
}

// ListByTag returns all skills with the given tag
func (m *Manager) ListByTag(tag string) ([]core.Skill, error) {
	allSkills, err := m.List()
	if err != nil {
		return nil, err
	}

	var filtered []core.Skill
	for _, skill := range allSkills {
		for _, t := range skill.Tags {
			if strings.EqualFold(t, tag) {
				filtered = append(filtered, skill)
				break
			}
		}
	}

	return filtered, nil
}

// ListByGroup returns all skills in the given group
func (m *Manager) ListByGroup(group string) ([]core.Skill, error) {
	allSkills, err := m.List()
	if err != nil {
		return nil, err
	}

	var filtered []core.Skill
	for _, skill := range allSkills {
		if strings.EqualFold(skill.Group, group) {
			filtered = append(filtered, skill)
		}
	}

	return filtered, nil
}

// GetAllTags returns all unique tags used across skills
func (m *Manager) GetAllTags() ([]string, error) {
	skills, err := m.List()
	if err != nil {
		return nil, err
	}

	tagSet := make(map[string]bool)
	for _, skill := range skills {
		for _, tag := range skill.Tags {
			tagSet[strings.ToLower(tag)] = true
		}
	}

	var tags []string
	for tag := range tagSet {
		tags = append(tags, tag)
	}
	sort.Strings(tags)

	return tags, nil
}

// GetAllGroups returns all unique groups used across skills
func (m *Manager) GetAllGroups() ([]string, error) {
	skills, err := m.List()
	if err != nil {
		return nil, err
	}

	groupSet := make(map[string]bool)
	for _, skill := range skills {
		if skill.Group != "" {
			groupSet[skill.Group] = true
		}
	}

	var groups []string
	for group := range groupSet {
		groups = append(groups, group)
	}
	sort.Strings(groups)

	return groups, nil
}

// Get retrieves a skill by name
func (m *Manager) Get(name string) (*core.Skill, error) {
	return m.load(name)
}

// Edit updates an existing skill
func (m *Manager) Edit(name string, skill *core.Skill) error {
	dir := m.skillDir(name)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("skill '%s' not found", name)
	}

	// If name changed, rename directory
	if name != skill.Name {
		newDir := m.skillDir(skill.Name)
		if err := os.Rename(dir, newDir); err != nil {
			return fmt.Errorf("failed to rename skill: %w", err)
		}
	}

	return m.save(skill)
}

// Remove deletes a skill
func (m *Manager) Remove(name string) error {
	dir := m.skillDir(name)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("skill '%s' not found", name)
	}
	return os.RemoveAll(dir)
}

// save writes a skill to SKILL.md with YAML frontmatter
func (m *Manager) save(skill *core.Skill) error {
	path := m.skillPath(skill.Name)

	// Build the SKILL.md content
	var content strings.Builder

	// YAML frontmatter with all fields
	content.WriteString("---\n")
	frontmatter := struct {
		Name         string            `yaml:"name"`
		Description  string            `yaml:"description"`
		Extends      string            `yaml:"extends,omitempty"`
		Includes     []string          `yaml:"includes,omitempty"`
		Tags         []string          `yaml:"tags,omitempty"`
		Group        string            `yaml:"group,omitempty"`
		Template     string            `yaml:"template,omitempty"`
		Variables    map[string]string `yaml:"variables,omitempty"`
		Author       string            `yaml:"author,omitempty"`
		Version      string            `yaml:"version,omitempty"`
		OutputFormat string            `yaml:"output_format,omitempty"`
		Context      *core.ContextConfig `yaml:"context,omitempty"`
		Hooks        *core.HooksConfig   `yaml:"hooks,omitempty"`
		Chain        []string          `yaml:"chain,omitempty"`
	}{
		Name:         skill.Name,
		Description:  skill.Description,
		Extends:      skill.Extends,
		Includes:     skill.Includes,
		Tags:         skill.Tags,
		Group:        skill.Group,
		Template:     skill.Template,
		Variables:    skill.Variables,
		Author:       skill.Author,
		Version:      skill.Version,
		OutputFormat: skill.OutputFormat,
		Context:      skill.Context,
		Hooks:        skill.Hooks,
		Chain:        skill.Chain,
	}
	fm, err := yaml.Marshal(frontmatter)
	if err != nil {
		return err
	}
	content.Write(fm)
	content.WriteString("---\n\n")

	// Markdown content with rules
	content.WriteString(fmt.Sprintf("# %s\n\n", skill.Name))
	content.WriteString(fmt.Sprintf("%s\n\n", skill.Description))

	if len(skill.Rules) > 0 {
		content.WriteString("## Rules\n\n")
		for _, rule := range skill.Rules {
			content.WriteString(fmt.Sprintf("- %s\n", rule))
		}
	}

	return os.WriteFile(path, []byte(content.String()), 0644)
}

// load reads a skill from SKILL.md with YAML frontmatter
func (m *Manager) load(name string) (*core.Skill, error) {
	path := m.skillPath(name)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	content := string(data)

	// Parse YAML frontmatter
	if !strings.HasPrefix(content, "---\n") {
		return nil, fmt.Errorf("invalid SKILL.md: missing frontmatter")
	}

	// Find the end of frontmatter
	endIdx := strings.Index(content[4:], "\n---")
	if endIdx == -1 {
		return nil, fmt.Errorf("invalid SKILL.md: unclosed frontmatter")
	}

	frontmatterStr := content[4 : 4+endIdx]
	markdownContent := content[4+endIdx+4:] // Skip past the closing ---

	// Parse frontmatter with all fields
	var frontmatter struct {
		Name         string            `yaml:"name"`
		Description  string            `yaml:"description"`
		Extends      string            `yaml:"extends,omitempty"`
		Includes     []string          `yaml:"includes,omitempty"`
		Tags         []string          `yaml:"tags,omitempty"`
		Group        string            `yaml:"group,omitempty"`
		Template     string            `yaml:"template,omitempty"`
		Variables    map[string]string `yaml:"variables,omitempty"`
		Author       string            `yaml:"author,omitempty"`
		Version      string            `yaml:"version,omitempty"`
		OutputFormat string            `yaml:"output_format,omitempty"`
		Context      *core.ContextConfig `yaml:"context,omitempty"`
		Hooks        *core.HooksConfig   `yaml:"hooks,omitempty"`
		Chain        []string          `yaml:"chain,omitempty"`
	}
	if err := yaml.Unmarshal([]byte(frontmatterStr), &frontmatter); err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	// Parse rules from markdown content
	rules := parseRulesFromMarkdown(markdownContent)

	return &core.Skill{
		Name:         frontmatter.Name,
		Description:  frontmatter.Description,
		Rules:        rules,
		Extends:      frontmatter.Extends,
		Includes:     frontmatter.Includes,
		Tags:         frontmatter.Tags,
		Group:        frontmatter.Group,
		Template:     frontmatter.Template,
		Variables:    frontmatter.Variables,
		Author:       frontmatter.Author,
		Version:      frontmatter.Version,
		OutputFormat: frontmatter.OutputFormat,
		Context:      frontmatter.Context,
		Hooks:        frontmatter.Hooks,
		Chain:        frontmatter.Chain,
	}, nil
}

// parseRulesFromMarkdown extracts rules from the markdown content
func parseRulesFromMarkdown(content string) []string {
	var rules []string
	inRulesSection := false

	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()

		// Check for Rules section header
		if strings.HasPrefix(line, "## Rules") {
			inRulesSection = true
			continue
		}

		// Check for next section (ends Rules section)
		if inRulesSection && strings.HasPrefix(line, "## ") {
			break
		}

		// Parse rule items
		if inRulesSection && strings.HasPrefix(line, "- ") {
			rule := strings.TrimPrefix(line, "- ")
			if rule != "" {
				rules = append(rules, rule)
			}
		}
	}

	return rules
}

// GetSkillDir returns the directory path for a skill (for history/rollback)
func (m *Manager) GetSkillDir(name string) string {
	return m.skillDir(name)
}

// ============== Version History ==============

// SaveVersion saves the current skill to history
func (m *Manager) SaveVersion(name string) error {
	skill, err := m.Get(name)
	if err != nil {
		return err
	}

	historyPath := filepath.Join(HistoryDir, strings.ToLower(name))
	if err := os.MkdirAll(historyPath, 0755); err != nil {
		return err
	}

	// Get next version number
	version := m.getNextVersion(name)
	versionFile := filepath.Join(historyPath, fmt.Sprintf("SKILL.v%d.md", version))

	// Read current file and copy to history
	currentPath := m.skillPath(name)
	data, err := os.ReadFile(currentPath)
	if err != nil {
		return err
	}

	// Add timestamp comment
	timestamp := time.Now().Format(time.RFC3339)
	header := fmt.Sprintf("<!-- Version %d saved at %s -->\n", version, timestamp)

	_ = skill // Used for potential future enhancements
	return os.WriteFile(versionFile, append([]byte(header), data...), 0644)
}

// GetVersions returns all versions for a skill
func (m *Manager) GetVersions(name string) ([]VersionInfo, error) {
	historyPath := filepath.Join(HistoryDir, strings.ToLower(name))
	if _, err := os.Stat(historyPath); os.IsNotExist(err) {
		return []VersionInfo{}, nil
	}

	entries, err := os.ReadDir(historyPath)
	if err != nil {
		return nil, err
	}

	var versions []VersionInfo
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasPrefix(entry.Name(), "SKILL.v") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		// Extract version number
		var version int
		fmt.Sscanf(entry.Name(), "SKILL.v%d.md", &version)

		versions = append(versions, VersionInfo{
			Version:   version,
			Timestamp: info.ModTime(),
			Path:      filepath.Join(historyPath, entry.Name()),
		})
	}

	// Sort by version descending
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].Version > versions[j].Version
	})

	return versions, nil
}

// VersionInfo holds version metadata
type VersionInfo struct {
	Version   int
	Timestamp time.Time
	Path      string
}

// getNextVersion returns the next version number for a skill
func (m *Manager) getNextVersion(name string) int {
	versions, err := m.GetVersions(name)
	if err != nil || len(versions) == 0 {
		return 1
	}
	return versions[0].Version + 1
}

// Rollback restores a skill to a previous version
func (m *Manager) Rollback(name string, version int) error {
	historyPath := filepath.Join(HistoryDir, strings.ToLower(name))
	versionFile := filepath.Join(historyPath, fmt.Sprintf("SKILL.v%d.md", version))

	if _, err := os.Stat(versionFile); os.IsNotExist(err) {
		return fmt.Errorf("version %d not found for skill '%s'", version, name)
	}

	// Save current as new version first
	if err := m.SaveVersion(name); err != nil {
		return fmt.Errorf("failed to save current version: %w", err)
	}

	// Read version file (skip the timestamp comment)
	data, err := os.ReadFile(versionFile)
	if err != nil {
		return err
	}

	content := string(data)
	// Skip the timestamp comment line
	if strings.HasPrefix(content, "<!--") {
		if idx := strings.Index(content, "-->\n"); idx != -1 {
			content = content[idx+4:]
		}
	}

	// Write to current skill file
	return os.WriteFile(m.skillPath(name), []byte(content), 0644)
}

// Diff returns the difference between two versions
func (m *Manager) Diff(name string, v1, v2 int) (string, string, error) {
	var content1, content2 string

	if v1 == 0 {
		// Compare with current
		data, err := os.ReadFile(m.skillPath(name))
		if err != nil {
			return "", "", err
		}
		content1 = string(data)
	} else {
		historyPath := filepath.Join(HistoryDir, strings.ToLower(name))
		data, err := os.ReadFile(filepath.Join(historyPath, fmt.Sprintf("SKILL.v%d.md", v1)))
		if err != nil {
			return "", "", err
		}
		content1 = string(data)
	}

	if v2 == 0 {
		data, err := os.ReadFile(m.skillPath(name))
		if err != nil {
			return "", "", err
		}
		content2 = string(data)
	} else {
		historyPath := filepath.Join(HistoryDir, strings.ToLower(name))
		data, err := os.ReadFile(filepath.Join(historyPath, fmt.Sprintf("SKILL.v%d.md", v2)))
		if err != nil {
			return "", "", err
		}
		content2 = string(data)
	}

	return content1, content2, nil
}

// ============== Export/Import ==============

// Export exports a skill to the specified format
func (m *Manager) Export(name string, format string) (string, error) {
	skill, err := m.Get(name)
	if err != nil {
		return "", err
	}

	switch format {
	case "json":
		data, err := json.MarshalIndent(skill, "", "  ")
		if err != nil {
			return "", err
		}
		return string(data), nil
	case "yaml":
		data, err := yaml.Marshal(skill)
		if err != nil {
			return "", err
		}
		return string(data), nil
	case "markdown", "md":
		// Return the raw SKILL.md content
		data, err := os.ReadFile(m.skillPath(name))
		if err != nil {
			return "", err
		}
		return string(data), nil
	default:
		return "", fmt.Errorf("unsupported export format: %s", format)
	}
}

// Import imports a skill from content in the specified format
func (m *Manager) Import(content string, format string) (*core.Skill, error) {
	var skill core.Skill

	switch format {
	case "json":
		if err := json.Unmarshal([]byte(content), &skill); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %w", err)
		}
	case "yaml":
		if err := yaml.Unmarshal([]byte(content), &skill); err != nil {
			return nil, fmt.Errorf("failed to parse YAML: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported import format: %s", format)
	}

	return &skill, nil
}

// ============== Templates ==============

// GetBuiltinTemplates returns all built-in skill templates
func GetBuiltinTemplates() []core.SkillTemplate {
	return []core.SkillTemplate{
		{
			Name:        "code-review",
			Description: "Review code for quality, bugs, and best practices",
			Category:    "development",
			Skill: core.Skill{
				Name:        "code-review",
				Description: "Reviews code for quality issues, potential bugs, security vulnerabilities, and adherence to best practices. Provides actionable feedback with specific line references.",
				Tags:        []string{"code", "review", "quality"},
				Rules: []string{
					"Always cite specific line numbers when referencing code issues",
					"Categorize issues by severity: critical, warning, suggestion",
					"Check for common security vulnerabilities (SQL injection, XSS, etc.)",
					"Verify error handling is comprehensive and appropriate",
					"Ensure code follows the project's established patterns and conventions",
					"Look for performance issues like N+1 queries, unnecessary loops",
					"Check for proper resource cleanup (file handles, connections)",
					"Verify tests cover the critical paths of new code",
				},
			},
		},
		{
			Name:        "commit-message",
			Description: "Generate conventional commit messages",
			Category:    "git",
			Skill: core.Skill{
				Name:        "commit-message",
				Description: "Generates clear, conventional commit messages following the Conventional Commits specification. Analyzes staged changes to determine the appropriate type and scope.",
				Tags:        []string{"git", "commit", "automation"},
				Rules: []string{
					"Use conventional commit format: type(scope): description",
					"Valid types: feat, fix, docs, style, refactor, test, chore, perf, ci",
					"Keep the subject line under 72 characters",
					"Use imperative mood in the subject (Add, not Added)",
					"Include a body for complex changes explaining the why",
					"Reference issue numbers when applicable",
					"Group related changes into a single commit",
					"Never include generated files or dependencies in the diff analysis",
				},
			},
		},
		{
			Name:        "documentation",
			Description: "Write clear technical documentation",
			Category:    "docs",
			Skill: core.Skill{
				Name:        "documentation",
				Description: "Creates clear, comprehensive technical documentation. Explains concepts at the appropriate level for the target audience and includes practical examples.",
				Tags:        []string{"docs", "writing", "technical"},
				Rules: []string{
					"Start with a clear one-sentence summary of what this documents",
					"Include a quick-start example within the first 3 sections",
					"Use consistent heading hierarchy (h2 for sections, h3 for subsections)",
					"Provide code examples for every API or function documented",
					"Include both success and error cases in examples",
					"Link to related documentation rather than duplicating content",
					"Use tables for comparing options or listing parameters",
					"End with a troubleshooting or FAQ section for complex topics",
				},
			},
		},
		{
			Name:        "testing",
			Description: "Write comprehensive test suites",
			Category:    "development",
			Skill: core.Skill{
				Name:        "testing",
				Description: "Designs and implements comprehensive test suites. Covers unit tests, integration tests, and edge cases with clear assertions and good test isolation.",
				Tags:        []string{"testing", "quality", "automation"},
				Rules: []string{
					"Follow Arrange-Act-Assert (AAA) pattern in all tests",
					"Name tests descriptively: should_[expected]_when_[condition]",
					"Test one behavior per test function",
					"Use test fixtures for shared setup, avoid test interdependence",
					"Include edge cases: empty inputs, nulls, boundaries, errors",
					"Mock external dependencies, don't make real network calls",
					"Verify both positive and negative test cases",
					"Aim for behavior coverage, not just line coverage",
				},
			},
		},
		{
			Name:        "debugging",
			Description: "Systematic debugging and root cause analysis",
			Category:    "development",
			Skill: core.Skill{
				Name:        "debugging",
				Description: "Systematic approach to debugging issues. Uses scientific method to isolate problems, identify root causes, and verify fixes don't introduce regressions.",
				Tags:        []string{"debugging", "troubleshooting", "analysis"},
				Rules: []string{
					"Reproduce the issue before attempting any fix",
					"Gather evidence: logs, stack traces, error messages",
					"Form a hypothesis about the root cause before making changes",
					"Isolate variables by testing one change at a time",
					"Check for recent changes that correlate with issue onset",
					"Verify the fix actually resolves the issue, don't assume",
					"Document the root cause and fix for future reference",
					"Consider if similar issues exist elsewhere in the codebase",
				},
			},
		},
		{
			Name:        "api-design",
			Description: "Design RESTful APIs following best practices",
			Category:    "architecture",
			Skill: core.Skill{
				Name:        "api-design",
				Description: "Designs RESTful APIs with consistent patterns, proper HTTP semantics, clear error handling, and good developer experience.",
				Tags:        []string{"api", "rest", "design"},
				Rules: []string{
					"Use nouns for resources, verbs come from HTTP methods",
					"Return appropriate HTTP status codes (201 for create, 204 for delete)",
					"Use consistent error response format with code, message, and details",
					"Version APIs in the URL path (/v1/, /v2/)",
					"Support pagination for list endpoints with limit/offset or cursor",
					"Use JSON:API or similar spec for response envelope structure",
					"Document all endpoints with request/response examples",
					"Implement proper CORS headers for browser clients",
				},
			},
		},
		{
			Name:        "security-review",
			Description: "Review code for security vulnerabilities",
			Category:    "security",
			Skill: core.Skill{
				Name:        "security-review",
				Description: "Audits code for security vulnerabilities following OWASP guidelines. Identifies injection flaws, authentication issues, data exposure, and other common security problems.",
				Tags:        []string{"security", "audit", "owasp"},
				Rules: []string{
					"Check all user input is validated and sanitized",
					"Verify SQL queries use parameterized statements",
					"Ensure authentication tokens are not logged or exposed",
					"Check for proper authorization on all endpoints",
					"Verify sensitive data is encrypted at rest and in transit",
					"Look for hardcoded secrets, keys, or credentials",
					"Check dependencies for known vulnerabilities",
					"Verify proper HTTPS/TLS configuration",
				},
			},
		},
		{
			Name:        "refactoring",
			Description: "Improve code structure without changing behavior",
			Category:    "development",
			Skill: core.Skill{
				Name:        "refactoring",
				Description: "Improves code structure, readability, and maintainability while preserving existing behavior. Uses established refactoring patterns and ensures tests pass.",
				Tags:        []string{"refactoring", "clean-code", "maintenance"},
				Rules: []string{
					"Ensure comprehensive tests exist before refactoring",
					"Make one refactoring change at a time, verify tests pass",
					"Extract methods when functions exceed 20-30 lines",
					"Replace magic numbers with named constants",
					"Apply DRY only when duplication is true duplication",
					"Prefer composition over inheritance for flexibility",
					"Keep the refactoring scope focused, avoid feature creep",
					"Document the rationale for significant structural changes",
				},
			},
		},
	}
}

// ============== Workspace ==============

// LoadWorkspace loads the workspace configuration
func LoadWorkspace() (*core.Workspace, error) {
	if _, err := os.Stat(WorkspaceFile); os.IsNotExist(err) {
		return nil, nil // No workspace configured
	}

	data, err := os.ReadFile(WorkspaceFile)
	if err != nil {
		return nil, err
	}

	var workspace core.Workspace
	if err := yaml.Unmarshal(data, &workspace); err != nil {
		return nil, err
	}

	return &workspace, nil
}

// SaveWorkspace saves the workspace configuration
func SaveWorkspace(workspace *core.Workspace) error {
	if err := os.MkdirAll(filepath.Dir(WorkspaceFile), 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(workspace)
	if err != nil {
		return err
	}

	return os.WriteFile(WorkspaceFile, data, 0644)
}
