package skills

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"openskill/pkg/core"

	"gopkg.in/yaml.v3"
)

const SkillsDir = ".claude/skills"

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

	// YAML frontmatter
	content.WriteString("---\n")
	frontmatter := struct {
		Name        string   `yaml:"name"`
		Description string   `yaml:"description"`
		Extends     string   `yaml:"extends,omitempty"`
		Includes    []string `yaml:"includes,omitempty"`
	}{
		Name:        skill.Name,
		Description: skill.Description,
		Extends:     skill.Extends,
		Includes:    skill.Includes,
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

	// Parse frontmatter
	var frontmatter struct {
		Name        string   `yaml:"name"`
		Description string   `yaml:"description"`
		Extends     string   `yaml:"extends,omitempty"`
		Includes    []string `yaml:"includes,omitempty"`
	}
	if err := yaml.Unmarshal([]byte(frontmatterStr), &frontmatter); err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	// Parse rules from markdown content
	rules := parseRulesFromMarkdown(markdownContent)

	return &core.Skill{
		Name:        frontmatter.Name,
		Description: frontmatter.Description,
		Rules:       rules,
		Extends:     frontmatter.Extends,
		Includes:    frontmatter.Includes,
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
