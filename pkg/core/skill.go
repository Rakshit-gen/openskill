package core

// Skill represents a Claude skill definition
type Skill struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Rules       []string `yaml:"rules,omitempty"`
	Extends     string   `yaml:"extends,omitempty"`  // Name of parent skill to inherit from
	Includes    []string `yaml:"includes,omitempty"` // Names of skills to compose/merge
}
