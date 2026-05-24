package config

import "fmt"

// StepOutputConfig controls how individual step output is captured and stored.
type StepOutputConfig struct {
	Enabled       bool              // capture step stdout/stderr into named variables
	MaxBytes      int               // maximum bytes to capture per step (0 = unlimited)
	TrimSpace     bool              // trim leading/trailing whitespace from captured output
	StepOverrides map[string]string // step name -> variable name override
}

// DefaultStepOutputConfig returns a StepOutputConfig with sensible defaults.
func DefaultStepOutputConfig() StepOutputConfig {
	return StepOutputConfig{
		Enabled:       false,
		MaxBytes:      65536, // 64 KiB
		TrimSpace:     true,
		StepOverrides: map[string]string{},
	}
}

// ValidateStepOutput returns an error if the config is invalid.
func ValidateStepOutput(c StepOutputConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.MaxBytes < 0 {
		return fmt.Errorf("step_output: max_bytes must be >= 0, got %d", c.MaxBytes)
	}
	if c.MaxBytes > 10*1024*1024 {
		return fmt.Errorf("step_output: max_bytes exceeds maximum allowed value of 10485760, got %d", c.MaxBytes)
	}
	for step, varName := range c.StepOverrides {
		if step == "" {
			return fmt.Errorf("step_output: step name in overrides must not be empty")
		}
		if varName == "" {
			return fmt.Errorf("step_output: variable name for step %q must not be empty", step)
		}
	}
	return nil
}
