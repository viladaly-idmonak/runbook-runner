package config

import "errors"

// DryRunConfig controls whether the runner executes commands or only simulates them.
type DryRunConfig struct {
	// Enabled causes all steps to be printed but not executed.
	Enabled bool `yaml:"enabled"`

	// PrintCommands controls whether each command is printed to stdout during dry-run.
	PrintCommands bool `yaml:"print_commands"`

	// ExitOnFirstSkip causes the runner to stop after the first skipped step (dry-run only).
	ExitOnFirstSkip bool `yaml:"exit_on_first_skip"`
}

// DefaultDryRunConfig returns a DryRunConfig with safe, non-destructive defaults.
func DefaultDryRunConfig() DryRunConfig {
	return DryRunConfig{
		Enabled:         false,
		PrintCommands:   true,
		ExitOnFirstSkip: false,
	}
}

// ValidateDryRun validates the DryRunConfig fields.
// ExitOnFirstSkip is only meaningful when Enabled is true.
func ValidateDryRun(cfg DryRunConfig) error {
	if !cfg.Enabled && cfg.ExitOnFirstSkip {
		return errors.New("dry_run: exit_on_first_skip requires enabled to be true")
	}
	return nil
}
