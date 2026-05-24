package config

import "fmt"

// StepDependsConfig controls step dependency ordering enforcement.
// When enabled, a step will be skipped (or fail) if any of its named
// dependencies did not complete successfully.
type StepDependsConfig struct {
	// Enabled turns dependency checking on or off.
	Enabled bool `yaml:"enabled"`

	// FailOnUnmet causes the runner to return an error when a dependency
	// was not satisfied, rather than silently skipping the step.
	FailOnUnmet bool `yaml:"fail_on_unmet"`

	// Deps maps a step name to the list of step names it depends on.
	Deps map[string][]string `yaml:"deps"`
}

// DefaultStepDependsConfig returns a StepDependsConfig with safe defaults.
func DefaultStepDependsConfig() StepDependsConfig {
	return StepDependsConfig{
		Enabled:     false,
		FailOnUnmet: true,
		Deps:        map[string][]string{},
	}
}

// ValidateStepDepends checks the StepDependsConfig for consistency.
func ValidateStepDepends(c StepDependsConfig) error {
	if !c.Enabled {
		return nil
	}
	for step, deps := range c.Deps {
		if step == "" {
			return fmt.Errorf("step_depends: dep entry has empty step name")
		}
		seen := map[string]bool{}
		for _, d := range deps {
			if d == "" {
				return fmt.Errorf("step_depends: step %q has empty dependency name", step)
			}
			if d == step {
				return fmt.Errorf("step_depends: step %q depends on itself", step)
			}
			if seen[d] {
				return fmt.Errorf("step_depends: step %q has duplicate dependency %q", step, d)
			}
			seen[d] = true
		}
	}
	return nil
}
