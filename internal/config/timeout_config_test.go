package config

import (
	"testing"
	"time"
)

func TestDefaultTimeoutConfig_Values(t *testing.T) {
	tc := DefaultTimeoutConfig()

	if tc.DefaultTimeout != 30*time.Second {
		t.Errorf("expected DefaultTimeout 30s, got %s", tc.DefaultTimeout)
	}
	if tc.MaxTimeout != 10*time.Minute {
		t.Errorf("expected MaxTimeout 10m, got %s", tc.MaxTimeout)
	}
}

func TestValidateTimeout_Valid(t *testing.T) {
	tc := DefaultTimeoutConfig()
	if err := ValidateTimeout(tc); err != nil {
		t.Errorf("expected no error for valid config, got: %v", err)
	}
}

func TestValidateTimeout_ZeroDefaultTimeout(t *testing.T) {
	tc := TimeoutConfig{
		DefaultTimeout: 0,
		MaxTimeout:     5 * time.Minute,
	}
	if err := ValidateTimeout(tc); err != nil {
		t.Errorf("zero DefaultTimeout should be valid (no timeout), got: %v", err)
	}
}

func TestValidateTimeout_ZeroBothMeansNoLimit(t *testing.T) {
	tc := TimeoutConfig{DefaultTimeout: 0, MaxTimeout: 0}
	if err := ValidateTimeout(tc); err != nil {
		t.Errorf("both zero should be valid, got: %v", err)
	}
}

func TestValidateTimeout_NegativeDefaultTimeout(t *testing.T) {
	tc := TimeoutConfig{
		DefaultTimeout: -1 * time.Second,
		MaxTimeout:     5 * time.Minute,
	}
	if err := ValidateTimeout(tc); err == nil {
		t.Error("expected error for negative DefaultTimeout, got nil")
	}
}

func TestValidateTimeout_NegativeMaxTimeout(t *testing.T) {
	tc := TimeoutConfig{
		DefaultTimeout: 10 * time.Second,
		MaxTimeout:     -1 * time.Minute,
	}
	if err := ValidateTimeout(tc); err == nil {
		t.Error("expected error for negative MaxTimeout, got nil")
	}
}

func TestValidateTimeout_DefaultExceedsMax(t *testing.T) {
	tc := TimeoutConfig{
		DefaultTimeout: 20 * time.Minute,
		MaxTimeout:     5 * time.Minute,
	}
	if err := ValidateTimeout(tc); err == nil {
		t.Error("expected error when DefaultTimeout exceeds MaxTimeout, got nil")
	}
}

func TestValidateTimeout_DefaultEqualToMax(t *testing.T) {
	tc := TimeoutConfig{
		DefaultTimeout: 5 * time.Minute,
		MaxTimeout:     5 * time.Minute,
	}
	if err := ValidateTimeout(tc); err != nil {
		t.Errorf("DefaultTimeout equal to MaxTimeout should be valid, got: %v", err)
	}
}
