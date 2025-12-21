package core

// Skill represents a Claude skill definition
type Skill struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Rules       []string `yaml:"rules,omitempty"`
}
