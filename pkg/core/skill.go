package core

// Skill represents a Claude skill definition
type Skill struct {
	Name        string   `yaml:"name" json:"name"`
	Description string   `yaml:"description" json:"description"`
	Rules       []string `yaml:"rules,omitempty" json:"rules,omitempty"`
	Extends     string   `yaml:"extends,omitempty" json:"extends,omitempty"`   // Name of parent skill to inherit from
	Includes    []string `yaml:"includes,omitempty" json:"includes,omitempty"` // Names of skills to compose/merge

	// New features
	Tags       []string          `yaml:"tags,omitempty" json:"tags,omitempty"`             // Categories/tags for organization
	Group      string            `yaml:"group,omitempty" json:"group,omitempty"`           // Skill group/bundle name
	Template   string            `yaml:"template,omitempty" json:"template,omitempty"`     // Template this skill was created from
	Variables  map[string]string `yaml:"variables,omitempty" json:"variables,omitempty"`   // Configurable parameters
	Author     string            `yaml:"author,omitempty" json:"author,omitempty"`         // Skill author
	Version    string            `yaml:"version,omitempty" json:"version,omitempty"`       // Semantic version
	OutputFormat string          `yaml:"output_format,omitempty" json:"output_format,omitempty"` // Expected output format (markdown, json, code)

	// Context providers
	Context    *ContextConfig    `yaml:"context,omitempty" json:"context,omitempty"`       // Context gathering configuration

	// Hooks
	Hooks      *HooksConfig      `yaml:"hooks,omitempty" json:"hooks,omitempty"`           // Pre/post execution hooks

	// Chaining/Workflows
	Chain      []string          `yaml:"chain,omitempty" json:"chain,omitempty"`           // Skills to run in sequence
}

// ContextConfig defines how a skill gathers context
type ContextConfig struct {
	Files       []string `yaml:"files,omitempty" json:"files,omitempty"`             // Files to read
	Globs       []string `yaml:"globs,omitempty" json:"globs,omitempty"`             // Glob patterns to match
	Commands    []string `yaml:"commands,omitempty" json:"commands,omitempty"`       // Commands to execute
	URLs        []string `yaml:"urls,omitempty" json:"urls,omitempty"`               // URLs to fetch
	Environment []string `yaml:"environment,omitempty" json:"environment,omitempty"` // Env vars to include
}

// HooksConfig defines pre/post execution hooks
type HooksConfig struct {
	Pre  []string `yaml:"pre,omitempty" json:"pre,omitempty"`   // Commands to run before skill execution
	Post []string `yaml:"post,omitempty" json:"post,omitempty"` // Commands to run after skill execution
}

// SkillGroup represents a bundle of related skills
type SkillGroup struct {
	Name        string   `yaml:"name" json:"name"`
	Description string   `yaml:"description" json:"description"`
	Skills      []string `yaml:"skills" json:"skills"`             // Skill names in this group
	Tags        []string `yaml:"tags,omitempty" json:"tags,omitempty"`
}

// SkillTemplate represents a pre-built skill template
type SkillTemplate struct {
	Name        string            `yaml:"name" json:"name"`
	Description string            `yaml:"description" json:"description"`
	Category    string            `yaml:"category" json:"category"`
	Variables   map[string]string `yaml:"variables,omitempty" json:"variables,omitempty"` // Default variable values
	Skill       Skill             `yaml:"skill" json:"skill"`                             // The template skill
}

// Workspace represents project-specific skill configuration
type Workspace struct {
	Name        string   `yaml:"name" json:"name"`
	Description string   `yaml:"description,omitempty" json:"description,omitempty"`
	Skills      []string `yaml:"skills,omitempty" json:"skills,omitempty"`   // Skills enabled in this workspace
	Groups      []string `yaml:"groups,omitempty" json:"groups,omitempty"`   // Groups enabled in this workspace
	Overrides   map[string]map[string]string `yaml:"overrides,omitempty" json:"overrides,omitempty"` // Variable overrides per skill
}
