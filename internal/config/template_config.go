package config

import (
	"errors"
	"regexp"
)

// TemplateConfig controls variable substitution in runbook step commands.
type TemplateConfig struct {
	// Enabled turns template rendering on or off.
	Enabled bool `yaml:"enabled"`

	// Delimiters defines the left and right delimiters for template variables.
	// Defaults to "{{" and "}}".
	LeftDelim  string `yaml:"left_delim"`
	RightDelim string `yaml:"right_delim"`

	// Vars holds static key-value pairs injected into every step command.
	Vars map[string]string `yaml:"vars"`

	// StrictMode causes rendering to fail if a referenced variable is undefined.
	StrictMode bool `yaml:"strict_mode"`
}

// DefaultTemplateConfig returns a TemplateConfig with sensible defaults.
func DefaultTemplateConfig() TemplateConfig {
	return TemplateConfig{
		Enabled:    false,
		LeftDelim:  "{{",
		RightDelim: "}}",
		Vars:       map[string]string{},
		StrictMode: false,
	}
}

var delimRe = regexp.MustCompile(`^\S+$`)

// ValidateTemplate returns an error if cfg contains invalid values.
func ValidateTemplate(cfg TemplateConfig) error {
	if !cfg.Enabled {
		return nil
	}
	if cfg.LeftDelim == "" {
		return errors.New("template: left_delim must not be empty")
	}
	if cfg.RightDelim == "" {
		return errors.New("template: right_delim must not be empty")
	}
	if cfg.LeftDelim == cfg.RightDelim {
		return errors.New("template: left_delim and right_delim must differ")
	}
	if !delimRe.MatchString(cfg.LeftDelim) {
		return errors.New("template: left_delim must not contain whitespace")
	}
	if !delimRe.MatchString(cfg.RightDelim) {
		return errors.New("template: right_delim must not contain whitespace")
	}
	for k := range cfg.Vars {
		if k == "" {
			return errors.New("template: vars key must not be empty")
		}
	}
	return nil
}
