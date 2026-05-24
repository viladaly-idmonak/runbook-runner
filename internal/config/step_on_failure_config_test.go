package config

import (
	"testing"
)

func TestDefaultStepOnFailureConfig_Values(t *testing.T) {
	c := DefaultStepOnFailureConfig()
	if c.Enabled {
		t.Error("expected Enabled=false")
	}
	if c.Default != OnFailureStop {
		t.Errorf("expected Default=stop, got %q", c.Default)
	}
	if c.Overrides == nil {
		t.Error("expected non-nil Overrides map")
	}
}

func TestValidateStepOnFailure_DisabledSkipsValidation(t *testing.T) {
	c := StepOnFailureConfig{
		Enabled: false,
		Default: "bad-action",
	}
	if err := ValidateStepOnFailure(c); err != nil {
		t.Errorf("expected no error when disabled, got %v", err)
	}
}

func TestValidateStepOnFailure_ValidConfig(t *testing.T) {
	c := StepOnFailureConfig{
		Enabled: true,
		Default: OnFailureContinue,
		Overrides: map[string]OnFailureAction{
			"deploy": OnFailureRollback,
		},
	}
	if err := ValidateStepOnFailure(c); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateStepOnFailure_UnknownDefault(t *testing.T) {
	c := StepOnFailureConfig{
		Enabled: true,
		Default: "explode",
	}
	if err := ValidateStepOnFailure(c); err == nil {
		t.Error("expected error for unknown default action")
	}
}

func TestValidateStepOnFailure_EmptyStepName(t *testing.T) {
	c := StepOnFailureConfig{
		Enabled: true,
		Default: OnFailureStop,
		Overrides: map[string]OnFailureAction{
			"": OnFailureContinue,
		},
	}
	if err := ValidateStepOnFailure(c); err == nil {
		t.Error("expected error for empty step name in overrides")
	}
}

func TestValidateStepOnFailure_UnknownOverrideAction(t *testing.T) {
	c := StepOnFailureConfig{
		Enabled: true,
		Default: OnFailureStop,
		Overrides: map[string]OnFailureAction{
			"build": "retry-forever",
		},
	}
	if err := ValidateStepOnFailure(c); err == nil {
		t.Error("expected error for unknown override action")
	}
}

func TestActionForStep_Disabled_ReturnsStop(t *testing.T) {
	c := StepOnFailureConfig{Enabled: false, Default: OnFailureContinue}
	if got := ActionForStep(c, "any"); got != OnFailureStop {
		t.Errorf("expected stop when disabled, got %q", got)
	}
}

func TestActionForStep_UsesOverride(t *testing.T) {
	c := StepOnFailureConfig{
		Enabled: true,
		Default: OnFailureStop,
		Overrides: map[string]OnFailureAction{"deploy": OnFailureRollback},
	}
	if got := ActionForStep(c, "deploy"); got != OnFailureRollback {
		t.Errorf("expected rollback, got %q", got)
	}
}

func TestActionForStep_FallsBackToDefault(t *testing.T) {
	c := StepOnFailureConfig{
		Enabled:   true,
		Default:   OnFailureContinue,
		Overrides: map[string]OnFailureAction{},
	}
	if got := ActionForStep(c, "unknown-step"); got != OnFailureContinue {
		t.Errorf("expected continue, got %q", got)
	}
}
