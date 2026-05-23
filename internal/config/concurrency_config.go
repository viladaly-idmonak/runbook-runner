package config

import (
	"errors"
	"fmt"
)

const (
	// DefaultWorkers is the number of parallel step workers when not specified.
	DefaultWorkers = 1
	// MaxWorkers is the upper bound accepted by ValidateConcurrency.
	MaxWorkers = 64
)

// ConcurrencyConfig controls how many runbook steps may execute in parallel.
type ConcurrencyConfig struct {
	// Workers is the maximum number of goroutines used to run steps.
	Workers int `yaml:"workers"`
	// FailFast stops scheduling new steps after the first failure when true.
	FailFast bool `yaml:"fail_fast"`
}

// DefaultConcurrencyConfig returns a safe, sequential configuration.
func DefaultConcurrencyConfig() ConcurrencyConfig {
	return ConcurrencyConfig{
		Workers:  DefaultWorkers,
		FailFast: true,
	}
}

// ValidateConcurrency returns an error if cfg contains invalid values.
func ValidateConcurrency(cfg ConcurrencyConfig) error {
	if cfg.Workers < 1 {
		return errors.New("concurrency: workers must be at least 1")
	}
	if cfg.Workers > MaxWorkers {
		return fmt.Errorf("concurrency: workers must not exceed %d", MaxWorkers)
	}
	return nil
}
