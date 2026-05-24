package config

import "fmt"

// ConditionMode controls how a step condition is evaluated.
type ConditionMode string

const (
	ConditionModeShell  ConditionMode = "shell"
	ConditionModeEnvSet ConditionMode = "env_set"
)

// ConditionConfig holds settings for conditional step execution.
type ConditionConfig struct {
	// Enabled turns conditional execution on or off.
	Enabled bool `yaml:"enabled"`

	// DefaultMode is the fallback evaluation mode when a step doesn't specify one.
	DefaultMode ConditionMode `yaml:"default_mode"`

	// SkipOnConditionFailure treats a failed condition as a skip rather than an error.
	SkipOnConditionFailure bool `yaml:"skip_on_condition_failure"`
}

// DefaultConditionConfig returns a ConditionConfig with sensible defaults.
func DefaultConditionConfig() ConditionConfig {
	return ConditionConfig{
		Enabled:                false,
		DefaultMode:            ConditionModeShell,
		SkipOnConditionFailure: true,
	}
}

// ValidateCondition returns an error if the config is invalid.
func ValidateCondition(c ConditionConfig) error {
	if !c.Enabled {
		return nil
	}
	switch c.DefaultMode {
	case ConditionModeShell, ConditionModeEnvSet:
		// valid
	case "":
		return fmt.Errorf("condition: default_mode must not be empty when enabled")
	default:
		return fmt.Errorf("condition: unknown default_mode %q (want \"shell\" or \"env_set\")", c.DefaultMode)
	}
	return nil
}
