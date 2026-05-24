package config

import "fmt"

// PromptConfig controls interactive confirmation prompts before executing steps.
type PromptConfig struct {
	// Enabled turns on interactive prompts before each step.
	Enabled bool `yaml:"enabled"`

	// OnFailure prompts the user before executing a rollback command.
	OnFailure bool `yaml:"on_failure"`

	// NonInteractive suppresses prompts and auto-confirms when true.
	// Useful for CI environments.
	NonInteractive bool `yaml:"non_interactive"`
}

// DefaultPromptConfig returns a PromptConfig with safe defaults.
func DefaultPromptConfig() PromptConfig {
	return PromptConfig{
		Enabled:        false,
		OnFailure:      false,
		NonInteractive: false,
	}
}

// ValidatePrompt returns an error if the PromptConfig is inconsistent.
func ValidatePrompt(c PromptConfig) error {
	if c.NonInteractive && c.Enabled {
		return fmt.Errorf("prompt: non_interactive cannot be combined with enabled=true")
	}
	if c.NonInteractive && c.OnFailure {
		return fmt.Errorf("prompt: non_interactive cannot be combined with on_failure=true")
	}
	return nil
}
