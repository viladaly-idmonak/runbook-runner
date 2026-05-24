package config

import (
	"fmt"
	"time"
)

// RateConfig controls step execution rate limiting.
type RateConfig struct {
	Enabled       bool          `yaml:"enabled"`
	MaxPerMinute  int           `yaml:"max_per_minute"`
	Burst         int           `yaml:"burst"`
	MinStepDelay  time.Duration `yaml:"min_step_delay"`
}

// DefaultRateConfig returns a RateConfig with sensible defaults (disabled).
func DefaultRateConfig() RateConfig {
	return RateConfig{
		Enabled:      false,
		MaxPerMinute: 60,
		Burst:        5,
		MinStepDelay: 0,
	}
}

// ValidateRate returns an error if the RateConfig is invalid.
func ValidateRate(r RateConfig) error {
	if !r.Enabled {
		return nil
	}
	if r.MaxPerMinute <= 0 {
		return fmt.Errorf("rate: max_per_minute must be greater than zero, got %d", r.MaxPerMinute)
	}
	if r.MaxPerMinute > 3600 {
		return fmt.Errorf("rate: max_per_minute must not exceed 3600, got %d", r.MaxPerMinute)
	}
	if r.Burst < 0 {
		return fmt.Errorf("rate: burst must be non-negative, got %d", r.Burst)
	}
	if r.Burst > r.MaxPerMinute {
		return fmt.Errorf("rate: burst (%d) must not exceed max_per_minute (%d)", r.Burst, r.MaxPerMinute)
	}
	if r.MinStepDelay < 0 {
		return fmt.Errorf("rate: min_step_delay must be non-negative, got %v", r.MinStepDelay)
	}
	return nil
}
