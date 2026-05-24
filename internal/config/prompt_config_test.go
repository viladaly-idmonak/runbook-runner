package config_test

import (
	"testing"

	"github.com/your-org/runbook-runner/internal/config"
)

func TestDefaultPromptConfig_Values(t *testing.T) {
	c := config.DefaultPromptConfig()
	if c.Enabled {
		t.Error("expected Enabled to be false")
	}
	if c.OnFailure {
		t.Error("expected OnFailure to be false")
	}
	if c.NonInteractive {
		t.Error("expected NonInteractive to be false")
	}
}

func TestValidatePrompt_Valid(t *testing.T) {
	c := config.PromptConfig{Enabled: true, OnFailure: true, NonInteractive: false}
	if err := config.ValidatePrompt(c); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidatePrompt_NonInteractiveWithEnabled(t *testing.T) {
	c := config.PromptConfig{Enabled: true, NonInteractive: true}
	if err := config.ValidatePrompt(c); err == nil {
		t.Error("expected error for non_interactive + enabled")
	}
}

func TestValidatePrompt_NonInteractiveWithOnFailure(t *testing.T) {
	c := config.PromptConfig{OnFailure: true, NonInteractive: true}
	if err := config.ValidatePrompt(c); err == nil {
		t.Error("expected error for non_interactive + on_failure")
	}
}

func TestValidatePrompt_AllDisabled(t *testing.T) {
	c := config.DefaultPromptConfig()
	if err := config.ValidatePrompt(c); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidatePrompt_NonInteractiveAlone(t *testing.T) {
	c := config.PromptConfig{NonInteractive: true}
	if err := config.ValidatePrompt(c); err != nil {
		t.Errorf("unexpected error for non_interactive alone: %v", err)
	}
}
