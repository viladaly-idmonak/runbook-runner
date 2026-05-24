package config

import "testing"

func TestDefaultStepEnvConfig_Values(t *testing.T) {
	c := DefaultStepEnvConfig()
	if c.Enabled {
		t.Error("expected Enabled=false")
	}
	if c.Overrides == nil {
		t.Error("expected non-nil Overrides map")
	}
}

func TestValidateStepEnv_DisabledSkipsValidation(t *testing.T) {
	c := StepEnvConfig{
		Enabled:   false,
		Overrides: map[string][]string{"step1": {""}}, // would fail if validated
	}
	if err := ValidateStepEnv(c); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestValidateStepEnv_Valid(t *testing.T) {
	c := StepEnvConfig{
		Enabled:   true,
		Overrides: map[string][]string{"deploy": {"REGION=us-east-1", "DEBUG=true"}},
	}
	if err := ValidateStepEnv(c); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateStepEnv_EmptyStepName(t *testing.T) {
	c := StepEnvConfig{
		Enabled:   true,
		Overrides: map[string][]string{"": {"KEY=val"}},
	}
	if err := ValidateStepEnv(c); err == nil {
		t.Error("expected error for empty step name")
	}
}

func TestValidateStepEnv_EmptyPair(t *testing.T) {
	c := StepEnvConfig{
		Enabled:   true,
		Overrides: map[string][]string{"step1": {""}},
	}
	if err := ValidateStepEnv(c); err == nil {
		t.Error("expected error for empty env pair")
	}
}

func TestValidateStepEnv_MissingEquals(t *testing.T) {
	c := StepEnvConfig{
		Enabled:   true,
		Overrides: map[string][]string{"step1": {"NOEQUALS"}},
	}
	if err := ValidateStepEnv(c); err == nil {
		t.Error("expected error for pair missing '='")
	}
}

func TestValidateStepEnv_EmptyKey(t *testing.T) {
	c := StepEnvConfig{
		Enabled:   true,
		Overrides: map[string][]string{"step1": {"=value"}},
	}
	if err := ValidateStepEnv(c); err == nil {
		t.Error("expected error for empty key portion")
	}
}
