package config

import (
	"encoding/json"
	"os"
	"time"
)

// Config holds runtime configuration for the runbook runner.
type Config struct {
	// DryRun disables actual command execution.
	DryRun bool `json:"dry_run"`

	// StopOnFailure halts execution when a step fails (default: true).
	StopOnFailure bool `json:"stop_on_failure"`

	// AutoRollback triggers rollback commands on step failure.
	AutoRollback bool `json:"auto_rollback"`

	// StepTimeout is the maximum duration allowed per step.
	StepTimeout time.Duration `json:"step_timeout"`

	// ShellPath is the shell used to execute commands.
	ShellPath string `json:"shell_path"`

	// OutputFormat controls reporter output ("text" or "json").
	OutputFormat string `json:"output_format"`
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		DryRun:        false,
		StopOnFailure: true,
		AutoRollback:  false,
		StepTimeout:   30 * time.Second,
		ShellPath:     "/bin/sh",
		OutputFormat:  "text",
	}
}

// LoadFile reads a JSON config file from path and merges it over defaults.
func LoadFile(path string) (*Config, error) {
	cfg := Default()

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Validate returns an error if the config contains invalid values.
func (c *Config) Validate() error {
	if c.ShellPath == "" {
		return ErrEmptyShellPath
	}
	if c.OutputFormat != "text" && c.OutputFormat != "json" {
		return ErrInvalidOutputFormat
	}
	if c.StepTimeout <= 0 {
		return ErrInvalidTimeout
	}
	return nil
}
