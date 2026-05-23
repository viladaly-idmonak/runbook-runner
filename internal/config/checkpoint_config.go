package config

import "fmt"

// CheckpointConfig controls step-level checkpointing so that a runbook
// can resume from the last successful step after a failure.
type CheckpointConfig struct {
	// Enabled turns checkpointing on or off.
	Enabled bool `yaml:"enabled"`

	// Dir is the directory where checkpoint state files are written.
	// Defaults to ".runbook-checkpoints".
	Dir string `yaml:"dir"`

	// ResumeOnRestart causes the runner to skip already-completed steps
	// when a checkpoint file is present.
	ResumeOnRestart bool `yaml:"resume_on_restart"`
}

// DefaultCheckpointConfig returns a CheckpointConfig with sensible defaults.
func DefaultCheckpointConfig() CheckpointConfig {
	return CheckpointConfig{
		Enabled:         false,
		Dir:             ".runbook-checkpoints",
		ResumeOnRestart: true,
	}
}

// ValidateCheckpoint returns an error if cfg contains invalid values.
func ValidateCheckpoint(cfg CheckpointConfig) error {
	if !cfg.Enabled {
		return nil
	}
	if cfg.Dir == "" {
		return fmt.Errorf("checkpoint: dir must not be empty when checkpointing is enabled")
	}
	return nil
}
