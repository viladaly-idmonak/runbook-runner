package config

import "testing"

func TestDefaultStepRetryConfig_Values(t *testing.T) {
	c := DefaultStepRetryConfig()
	if c.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if c.Overrides == nil {
		t.Error("expected Overrides map to be initialised")
	}
}

func TestValidateStepRetry_DisabledSkipsValidation(t *testing.T) {
	c := StepRetryConfig{
		Enabled:   false,
		Overrides: map[string]RetryOverride{"step1": {MaxAttempts: 0, DelayMs: -1}},
	}
	if err := ValidateStepRetry(c); err != nil {
		t.Errorf("expected no error when disabled, got %v", err)
	}
}

func TestValidateStepRetry_ValidOverride(t *testing.T) {
	c := StepRetryConfig{
		Enabled:   true,
		Overrides: map[string]RetryOverride{"deploy": {MaxAttempts: 3, DelayMs: 500}},
	}
	if err := ValidateStepRetry(c); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateStepRetry_ZeroMaxAttempts(t *testing.T) {
	c := StepRetryConfig{
		Enabled:   true,
		Overrides: map[string]RetryOverride{"step1": {MaxAttempts: 0, DelayMs: 0}},
	}
	if err := ValidateStepRetry(c); err == nil {
		t.Error("expected error for max_attempts=0")
	}
}

func TestValidateStepRetry_TooManyAttempts(t *testing.T) {
	c := StepRetryConfig{
		Enabled:   true,
		Overrides: map[string]RetryOverride{"step1": {MaxAttempts: 21, DelayMs: 0}},
	}
	if err := ValidateStepRetry(c); err == nil {
		t.Error("expected error for max_attempts=21")
	}
}

func TestValidateStepRetry_NegativeDelay(t *testing.T) {
	c := StepRetryConfig{
		Enabled:   true,
		Overrides: map[string]RetryOverride{"step1": {MaxAttempts: 2, DelayMs: -10}},
	}
	if err := ValidateStepRetry(c); err == nil {
		t.Error("expected error for negative delay_ms")
	}
}

func TestLookupOverride_ReturnsDefault_WhenDisabled(t *testing.T) {
	c := StepRetryConfig{
		Enabled:   false,
		Overrides: map[string]RetryOverride{"step1": {MaxAttempts: 5, DelayMs: 100}},
	}
	def := RetryOverride{MaxAttempts: 1, DelayMs: 0}
	got := c.LookupOverride("step1", def)
	if got.MaxAttempts != 1 {
		t.Errorf("expected default MaxAttempts=1, got %d", got.MaxAttempts)
	}
}

func TestLookupOverride_ReturnsOverride_WhenEnabled(t *testing.T) {
	c := StepRetryConfig{
		Enabled:   true,
		Overrides: map[string]RetryOverride{"deploy": {MaxAttempts: 4, DelayMs: 200}},
	}
	def := RetryOverride{MaxAttempts: 1, DelayMs: 0}
	got := c.LookupOverride("deploy", def)
	if got.MaxAttempts != 4 {
		t.Errorf("expected override MaxAttempts=4, got %d", got.MaxAttempts)
	}
}

func TestLookupOverride_FallsBackToDefault_WhenNoMatch(t *testing.T) {
	c := StepRetryConfig{
		Enabled:   true,
		Overrides: map[string]RetryOverride{},
	}
	def := RetryOverride{MaxAttempts: 2, DelayMs: 50}
	got := c.LookupOverride("unknown", def)
	if got.MaxAttempts != 2 {
		t.Errorf("expected fallback MaxAttempts=2, got %d", got.MaxAttempts)
	}
}
