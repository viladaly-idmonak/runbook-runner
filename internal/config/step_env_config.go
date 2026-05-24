package config

import "fmt"

// StepEnvConfig controls per-step environment variable overrides.
type StepEnvConfig struct {
	// Enabled toggles per-step env injection.
	Enabled bool `yaml:"enabled"`

	// Overrides maps step names to key=value env pairs that are merged
	// into the process environment before that step runs.
	Overrides map[string][]string `yaml:"overrides"`
}

// DefaultStepEnvConfig returns a StepEnvConfig with safe defaults.
func DefaultStepEnvConfig() StepEnvConfig {
	return StepEnvConfig{
		Enabled:   false,
		Overrides: map[string][]string{},
	}
}

// ValidateStepEnv checks that every override entry is well-formed.
func ValidateStepEnv(c StepEnvConfig) error {
	if !c.Enabled {
		return nil
	}
	for stepName, pairs := range c.Overrides {
		if stepName == "" {
			return fmt.Errorf("step_env: override has empty step name")
		}
		for _, pair := range pairs {
			if pair == "" {
				return fmt.Errorf("step_env: step %q has empty env entry", stepName)
			}
			if !isValidEnvPair(pair) {
				return fmt.Errorf("step_env: step %q has invalid env pair %q (must be KEY=VALUE)", stepName, pair)
			}
		}
	}
	return nil
}

// isValidEnvPair returns true when s contains at least one '=' with a
// non-empty key portion.
func isValidEnvPair(s string) bool {
	for i, ch := range s {
		if ch == '=' {
			return i > 0
		}
	}
	return false
}
