package config_test

import (
	"testing"

	"github.com/example/runbook-runner/internal/config"
)

func TestDefaultConditionConfig_Values(t *testing.T) {
	c := config.DefaultConditionConfig()
	if c.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if c.DefaultMode != config.ConditionModeShell {
		t.Errorf("expected DefaultMode %q, got %q", config.ConditionModeShell, c.DefaultMode)
	}
	if !c.SkipOnConditionFailure {
		t.Error("expected SkipOnConditionFailure to be true by default")
	}
}

func TestValidateCondition_DisabledSkipsValidation(t *testing.T) {
	c := config.ConditionConfig{
		Enabled:     false,
		DefaultMode: "bad-value",
	}
	if err := config.ValidateCondition(c); err != nil {
		t.Errorf("expected no error when disabled, got %v", err)
	}
}

func TestValidateCondition_ValidShellMode(t *testing.T) {
	c := config.ConditionConfig{
		Enabled:     true,
		DefaultMode: config.ConditionModeShell,
	}
	if err := config.ValidateCondition(c); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateCondition_ValidEnvSetMode(t *testing.T) {
	c := config.ConditionConfig{
		Enabled:     true,
		DefaultMode: config.ConditionModeEnvSet,
	}
	if err := config.ValidateCondition(c); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateCondition_EmptyMode(t *testing.T) {
	c := config.ConditionConfig{
		Enabled:     true,
		DefaultMode: "",
	}
	if err := config.ValidateCondition(c); err == nil {
		t.Error("expected error for empty default_mode, got nil")
	}
}

func TestValidateCondition_UnknownMode(t *testing.T) {
	c := config.ConditionConfig{
		Enabled:     true,
		DefaultMode: "regex",
	}
	err := config.ValidateCondition(c)
	if err == nil {
		t.Fatal("expected error for unknown mode, got nil")
	}
	expected := "condition: unknown default_mode"
	if len(err.Error()) < len(expected) || err.Error()[:len(expected)] != expected {
		t.Errorf("unexpected error message: %v", err)
	}
}
