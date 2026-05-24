package config

import (
	"fmt"
	"path/filepath"
)

// StepWorkdirConfig controls per-step working directory overrides.
type StepWorkdirConfig struct {
	Enabled   bool              `yaml:"enabled"`
	Default   string            `yaml:"default"`
	Overrides map[string]string `yaml:"overrides"`
}

// DefaultStepWorkdirConfig returns a StepWorkdirConfig with safe defaults.
func DefaultStepWorkdirConfig() StepWorkdirConfig {
	return StepWorkdirConfig{
		Enabled:   false,
		Default:   "",
		Overrides: map[string]string{},
	}
}

// ValidateStepWorkdir checks that the config is consistent.
func ValidateStepWorkdir(c StepWorkdirConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.Default != "" && !filepath.IsAbs(c.Default) {
		return fmt.Errorf("step_workdir: default path %q must be absolute", c.Default)
	}
	for step, dir := range c.Overrides {
		if step == "" {
			return fmt.Errorf("step_workdir: override step name must not be empty")
		}
		if dir == "" {
			return fmt.Errorf("step_workdir: override directory for step %q must not be empty", step)
		}
		if !filepath.IsAbs(dir) {
			return fmt.Errorf("step_workdir: override path %q for step %q must be absolute", dir, step)
		}
	}
	return nil
}

// DirForStep returns the working directory to use for the named step.
// It returns the override if present, then the default, then "".
func (c StepWorkdirConfig) DirForStep(stepName string) string {
	if !c.Enabled {
		return ""
	}
	if dir, ok := c.Overrides[stepName]; ok {
		return dir
	}
	return c.Default
}
