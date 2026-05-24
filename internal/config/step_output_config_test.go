package config

import (
	"testing"
)

func TestDefaultStepOutputConfig_Values(t *testing.T) {
	c := DefaultStepOutputConfig()
	if c.Enabled {
		t.Error("expected Enabled to be false")
	}
	if c.MaxBytes != 65536 {
		t.Errorf("expected MaxBytes 65536, got %d", c.MaxBytes)
	}
	if !c.TrimSpace {
		t.Error("expected TrimSpace to be true")
	}
	if c.StepOverrides == nil {
		t.Error("expected StepOverrides to be non-nil")
	}
}

func TestValidateStepOutput_DisabledSkipsValidation(t *testing.T) {
	c := StepOutputConfig{Enabled: false, MaxBytes: -999}
	if err := ValidateStepOutput(c); err != nil {
		t.Errorf("expected no error when disabled, got %v", err)
	}
}

func TestValidateStepOutput_Valid(t *testing.T) {
	c := StepOutputConfig{
		Enabled:       true,
		MaxBytes:      1024,
		TrimSpace:     true,
		StepOverrides: map[string]string{"deploy": "DEPLOY_OUT"},
	}
	if err := ValidateStepOutput(c); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateStepOutput_NegativeMaxBytes(t *testing.T) {
	c := StepOutputConfig{Enabled: true, MaxBytes: -1, StepOverrides: map[string]string{}}
	if err := ValidateStepOutput(c); err == nil {
		t.Error("expected error for negative MaxBytes")
	}
}

func TestValidateStepOutput_ExceedsMaxBytes(t *testing.T) {
	c := StepOutputConfig{Enabled: true, MaxBytes: 20 * 1024 * 1024, StepOverrides: map[string]string{}}
	if err := ValidateStepOutput(c); err == nil {
		t.Error("expected error for MaxBytes exceeding limit")
	}
}

func TestValidateStepOutput_EmptyStepName(t *testing.T) {
	c := StepOutputConfig{
		Enabled:       true,
		MaxBytes:      1024,
		StepOverrides: map[string]string{"": "SOME_VAR"},
	}
	if err := ValidateStepOutput(c); err == nil {
		t.Error("expected error for empty step name in overrides")
	}
}

func TestValidateStepOutput_EmptyVarName(t *testing.T) {
	c := StepOutputConfig{
		Enabled:       true,
		MaxBytes:      1024,
		StepOverrides: map[string]string{"deploy": ""},
	}
	if err := ValidateStepOutput(c); err == nil {
		t.Error("expected error for empty variable name in overrides")
	}
}

func TestValidateStepOutput_ZeroMaxBytesAllowed(t *testing.T) {
	c := StepOutputConfig{Enabled: true, MaxBytes: 0, StepOverrides: map[string]string{}}
	if err := ValidateStepOutput(c); err != nil {
		t.Errorf("expected no error for zero MaxBytes (unlimited), got %v", err)
	}
}
