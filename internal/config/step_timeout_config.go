package config

import (
	"errors"
	"time"
)

// StepTimeoutConfig holds per-step timeout overrides keyed by step name.
type StepTimeoutConfig struct {
	// Enabled controls whether per-step timeout overrides are applied.
	Enabled bool `yaml:"enabled"`

	// Overrides maps a step name to a duration string (e.g. "30s", "2m").
	Overrides map[string]string `yaml:"overrides"`
}

// DefaultStepTimeoutConfig returns a StepTimeoutConfig with safe defaults.
func DefaultStepTimeoutConfig() StepTimeoutConfig {
	return StepTimeoutConfig{
		Enabled:   false,
		Overrides: map[string]string{},
	}
}

// ValidateStepTimeout checks that all override values are valid durations
// and that no step name is empty.
func ValidateStepTimeout(c StepTimeoutConfig) error {
	if !c.Enabled {
		return nil
	}
	if len(c.Overrides) == 0 {
		return errors.New("step_timeout: enabled but no overrides provided")
	}
	for name, raw := range c.Overrides {
		if name == "" {
			return errors.New("step_timeout: override key must not be empty")
		}
		d, err := time.ParseDuration(raw)
		if err != nil {
			return fmt.Errorf("step_timeout: invalid duration %q for step %q: %w", raw, name, err)
		}
		if d <= 0 {
			return fmt.Errorf("step_timeout: duration for step %q must be positive, got %s", name, raw)
		}
	}
	return nil
}

// Resolve returns the timeout for the given step name, falling back to
// defaultTimeout when no override exists or the config is disabled.
// A zero defaultTimeout means no limit.
func (c StepTimeoutConfig) Resolve(stepName string, defaultTimeout time.Duration) time.Duration {
	if !c.Enabled {
		return defaultTimeout
	}
	raw, ok := c.Overrides[stepName]
	if !ok {
		return defaultTimeout
	}
	d, err := time.ParseDuration(raw)
	if err != nil || d <= 0 {
		return defaultTimeout
	}
	return d
}
