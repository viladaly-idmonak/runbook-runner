package config

import "testing"

func TestDefaultStepLimitConfig_Values(t *testing.T) {
	c := DefaultStepLimitConfig()
	if c.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if c.MaxSteps != 100 {
		t.Errorf("expected MaxSteps=100, got %d", c.MaxSteps)
	}
	if c.FailOnExceed {
		t.Error("expected FailOnExceed to be false by default")
	}
}

func TestValidateStepLimit_DisabledSkipsValidation(t *testing.T) {
	c := StepLimitConfig{Enabled: false, MaxSteps: 0}
	if err := ValidateStepLimit(c); err != nil {
		t.Errorf("expected no error when disabled, got %v", err)
	}
}

func TestValidateStepLimit_ValidConfig(t *testing.T) {
	c := StepLimitConfig{Enabled: true, MaxSteps: 50}
	if err := ValidateStepLimit(c); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateStepLimit_ZeroMaxSteps(t *testing.T) {
	c := StepLimitConfig{Enabled: true, MaxSteps: 0}
	if err := ValidateStepLimit(c); err == nil {
		t.Error("expected error for MaxSteps=0")
	}
}

func TestValidateStepLimit_NegativeMaxSteps(t *testing.T) {
	c := StepLimitConfig{Enabled: true, MaxSteps: -5}
	if err := ValidateStepLimit(c); err == nil {
		t.Error("expected error for negative MaxSteps")
	}
}

func TestValidateStepLimit_ExceedsMaxAllowed(t *testing.T) {
	c := StepLimitConfig{Enabled: true, MaxSteps: 99_999}
	if err := ValidateStepLimit(c); err == nil {
		t.Error("expected error when MaxSteps exceeds 10000")
	}
}

func TestValidateStepLimit_BoundaryMaxSteps(t *testing.T) {
	c := StepLimitConfig{Enabled: true, MaxSteps: 10_000}
	if err := ValidateStepLimit(c); err != nil {
		t.Errorf("expected no error at boundary 10000, got %v", err)
	}
}
