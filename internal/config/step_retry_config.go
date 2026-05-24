package config

import "fmt"

// StepRetryConfig holds per-step retry overrides keyed by step name.
type StepRetryConfig struct {
	Enabled   bool                    `yaml:"enabled"`
	Overrides map[string]RetryOverride `yaml:"overrides"`
}

// RetryOverride defines retry parameters for a specific step.
type RetryOverride struct {
	MaxAttempts int `yaml:"max_attempts"`
	DelayMs     int `yaml:"delay_ms"`
}

// DefaultStepRetryConfig returns a StepRetryConfig with safe defaults.
func DefaultStepRetryConfig() StepRetryConfig {
	return StepRetryConfig{
		Enabled:   false,
		Overrides: map[string]RetryOverride{},
	}
}

// ValidateStepRetry checks that all per-step retry overrides are valid.
func ValidateStepRetry(c StepRetryConfig) error {
	if !c.Enabled {
		return nil
	}
	for name, o := range c.Overrides {
		if name == "" {
			return fmt.Errorf("step_retry: override key must not be empty")
		}
		if o.MaxAttempts < 1 {
			return fmt.Errorf("step_retry: override %q max_attempts must be >= 1, got %d", name, o.MaxAttempts)
		}
		if o.MaxAttempts > 20 {
			return fmt.Errorf("step_retry: override %q max_attempts must be <= 20, got %d", name, o.MaxAttempts)
		}
		if o.DelayMs < 0 {
			return fmt.Errorf("step_retry: override %q delay_ms must be >= 0, got %d", name, o.DelayMs)
		}
	}
	return nil
}

// LookupOverride returns the RetryOverride for a step name, or the provided
// default if no override is configured or the feature is disabled.
func (c StepRetryConfig) LookupOverride(stepName string, def RetryOverride) RetryOverride {
	if !c.Enabled {
		return def
	}
	if o, ok := c.Overrides[stepName]; ok {
		return o
	}
	return def
}
