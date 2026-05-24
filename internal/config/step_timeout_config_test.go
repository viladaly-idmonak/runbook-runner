package config

import (
	"testing"
	"time"
)

func TestDefaultStepTimeoutConfig_Values(t *testing.T) {
	c := DefaultStepTimeoutConfig()
	if c.Enabled {
		t.Error("expected Enabled to be false")
	}
	if len(c.Overrides) != 0 {
		t.Errorf("expected empty Overrides, got %v", c.Overrides)
	}
}

func TestValidateStepTimeout_DisabledSkipsValidation(t *testing.T) {
	c := StepTimeoutConfig{Enabled: false, Overrides: map[string]string{"step1": "INVALID"}}
	if err := ValidateStepTimeout(c); err != nil {
		t.Errorf("expected no error when disabled, got %v", err)
	}
}

func TestValidateStepTimeout_EnabledNoOverrides(t *testing.T) {
	c := StepTimeoutConfig{Enabled: true, Overrides: map[string]string{}}
	if err := ValidateStepTimeout(c); err == nil {
		t.Error("expected error for enabled config with no overrides")
	}
}

func TestValidateStepTimeout_ValidOverride(t *testing.T) {
	c := StepTimeoutConfig{
		Enabled:   true,
		Overrides: map[string]string{"deploy": "45s", "test": "2m"},
	}
	if err := ValidateStepTimeout(c); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateStepTimeout_InvalidDuration(t *testing.T) {
	c := StepTimeoutConfig{
		Enabled:   true,
		Overrides: map[string]string{"step1": "notaduration"},
	}
	if err := ValidateStepTimeout(c); err == nil {
		t.Error("expected error for invalid duration")
	}
}

func TestValidateStepTimeout_ZeroDuration(t *testing.T) {
	c := StepTimeoutConfig{
		Enabled:   true,
		Overrides: map[string]string{"step1": "0s"},
	}
	if err := ValidateStepTimeout(c); err == nil {
		t.Error("expected error for zero duration")
	}
}

func TestStepTimeoutConfig_Resolve_Disabled(t *testing.T) {
	c := DefaultStepTimeoutConfig()
	got := c.Resolve("step1", 10*time.Second)
	if got != 10*time.Second {
		t.Errorf("expected 10s, got %v", got)
	}
}

func TestStepTimeoutConfig_Resolve_Override(t *testing.T) {
	c := StepTimeoutConfig{
		Enabled:   true,
		Overrides: map[string]string{"deploy": "90s"},
	}
	got := c.Resolve("deploy", 30*time.Second)
	if got != 90*time.Second {
		t.Errorf("expected 90s, got %v", got)
	}
}

func TestStepTimeoutConfig_Resolve_FallbackWhenNoMatch(t *testing.T) {
	c := StepTimeoutConfig{
		Enabled:   true,
		Overrides: map[string]string{"other": "5s"},
	}
	got := c.Resolve("deploy", 20*time.Second)
	if got != 20*time.Second {
		t.Errorf("expected fallback 20s, got %v", got)
	}
}
