package config

import (
	"fmt"
	"time"
)

// RetryConfig holds retry settings for step execution.
type RetryConfig struct {
	MaxAttempts int           `yaml:"max_attempts"`
	Delay       time.Duration `yaml:"delay"`
}

// DefaultRetryConfig returns conservative retry defaults.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 1,
		Delay:       0,
	}
}

// ValidateRetry checks that the retry configuration is sensible.
func ValidateRetry(r RetryConfig) error {
	if r.MaxAttempts < 1 {
		return fmt.Errorf("%w: max_attempts must be >= 1, got %d", ErrInvalidConfig, r.MaxAttempts)
	}
	if r.MaxAttempts > 10 {
		return fmt.Errorf("%w: max_attempts must be <= 10, got %d", ErrInvalidConfig, r.MaxAttempts)
	}
	if r.Delay < 0 {
		return fmt.Errorf("%w: delay must be non-negative", ErrInvalidConfig)
	}
	if r.Delay > 5*time.Minute {
		return fmt.Errorf("%w: delay must be <= 5m, got %v", ErrInvalidConfig, r.Delay)
	}
	return nil
}
