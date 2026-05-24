package config

import "fmt"

// StepLimitConfig controls how many steps may be executed in a single run.
type StepLimitConfig struct {
	// Enabled turns step-limit enforcement on or off.
	Enabled bool `yaml:"enabled"`

	// MaxSteps is the maximum number of steps allowed to execute.
	// Must be >= 1 when Enabled is true.
	MaxSteps int `yaml:"max_steps"`

	// FailOnExceed causes the runner to return an error when the limit is
	// exceeded instead of silently stopping.
	FailOnExceed bool `yaml:"fail_on_exceed"`
}

// DefaultStepLimitConfig returns a safe default (disabled).
func DefaultStepLimitConfig() StepLimitConfig {
	return StepLimitConfig{
		Enabled:      false,
		MaxSteps:     100,
		FailOnExceed: false,
	}
}

// ValidateStepLimit returns an error when the configuration is inconsistent.
func ValidateStepLimit(c StepLimitConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.MaxSteps < 1 {
		return fmt.Errorf("step_limit: max_steps must be >= 1 when enabled, got %d", c.MaxSteps)
	}
	if c.MaxSteps > 10_000 {
		return fmt.Errorf("step_limit: max_steps exceeds maximum allowed value of 10000, got %d", c.MaxSteps)
	}
	return nil
}
