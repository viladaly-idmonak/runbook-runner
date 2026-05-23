package config

import "fmt"

// OutputFormat controls how step results are rendered.
type OutputFormat string

const (
	OutputFormatText OutputFormat = "text"
	OutputFormatJSON OutputFormat = "json"
)

// OutputConfig holds settings that control how run results are reported.
type OutputConfig struct {
	// Format selects the output renderer (text or json).
	Format OutputFormat `yaml:"format"`
	// Verbose enables per-step command output in the report.
	Verbose bool `yaml:"verbose"`
	// Color enables ANSI colour codes when Format is text.
	Color bool `yaml:"color"`
	// TimestampsEnabled adds wall-clock timestamps to each step line.
	TimestampsEnabled bool `yaml:"timestamps_enabled"`
}

// DefaultOutputConfig returns a safe, human-friendly default.
func DefaultOutputConfig() OutputConfig {
	return OutputConfig{
		Format:            OutputFormatText,
		Verbose:           false,
		Color:             true,
		TimestampsEnabled: false,
	}
}

// ValidateOutput returns an error when the config contains invalid values.
func ValidateOutput(c OutputConfig) error {
	switch c.Format {
	case OutputFormatText, OutputFormatJSON:
		// valid
	case "":
		return fmt.Errorf("output format must not be empty")
	default:
		return fmt.Errorf("unknown output format %q: must be \"text\" or \"json\"", c.Format)
	}
	return nil
}
