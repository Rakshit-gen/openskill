package skills

import (
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

// skillPath returns the file path for a skill
func (m *Manager) skillPath(name string) string {
	safeName := strings.ReplaceAll(strings.ToLower(name), " ", "-")
	return filepath.Join(m.baseDir, safeName+".yaml")
}

// Add creates a new skill
func (m *Manager) Add(skill *core.Skill) error {
	if err := m.ensureDir(); err != nil {
		return fmt.Errorf("failed to create skills directory: %w", err)
	}

	path := m.skillPath(skill.Name)
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("skill '%s' already exists", skill.Name)
	}

	return m.save(path, skill)
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
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}

		skill, err := m.load(filepath.Join(m.baseDir, entry.Name()))
		if err != nil {
			continue
		}
		skills = append(skills, *skill)
	}

	return skills, nil
}

// Get retrieves a skill by name
func (m *Manager) Get(name string) (*core.Skill, error) {
	path := m.skillPath(name)
	return m.load(path)
}

// Edit updates an existing skill
func (m *Manager) Edit(name string, skill *core.Skill) error {
	path := m.skillPath(name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("skill '%s' not found", name)
	}

	// If name changed, remove old file
	if name != skill.Name {
		os.Remove(path)
		path = m.skillPath(skill.Name)
	}

	return m.save(path, skill)
}

// Remove deletes a skill
func (m *Manager) Remove(name string) error {
	path := m.skillPath(name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("skill '%s' not found", name)
	}
	return os.Remove(path)
}

func (m *Manager) save(path string, skill *core.Skill) error {
	data, err := yaml.Marshal(skill)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (m *Manager) load(path string) (*core.Skill, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var skill core.Skill
	if err := yaml.Unmarshal(data, &skill); err != nil {
		return nil, err
	}
	return &skill, nil
}
