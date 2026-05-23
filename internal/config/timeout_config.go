package config

import (
	"fmt"
	"time"
)

// TimeoutConfig holds timeout-related configuration for runbook execution.
type TimeoutConfig struct {
	// DefaultTimeout is applied to every step unless overridden per-step.
	// A zero value means no timeout.
	DefaultTimeout time.Duration `yaml:"default_timeout"`

	// MaxTimeout is the upper bound for any per-step timeout override.
	// A zero value means no upper bound is enforced.
	MaxTimeout time.Duration `yaml:"max_timeout"`
}

// DefaultTimeoutConfig returns a TimeoutConfig with sensible defaults.
func DefaultTimeoutConfig() TimeoutConfig {
	return TimeoutConfig{
		DefaultTimeout: 30 * time.Second,
		MaxTimeout:     10 * time.Minute,
	}
}

// ValidateTimeout returns an error if the TimeoutConfig contains invalid values.
func ValidateTimeout(tc TimeoutConfig) error {
	if tc.DefaultTimeout < 0 {
		return fmt.Errorf("timeout: default_timeout must not be negative, got %s", tc.DefaultTimeout)
	}
	if tc.MaxTimeout < 0 {
		return fmt.Errorf("timeout: max_timeout must not be negative, got %s", tc.MaxTimeout)
	}
	if tc.MaxTimeout > 0 && tc.DefaultTimeout > tc.MaxTimeout {
		return fmt.Errorf(
			"timeout: default_timeout (%s) must not exceed max_timeout (%s)",
			tc.DefaultTimeout, tc.MaxTimeout,
		)
	}
	return nil
}
