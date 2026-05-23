package config

import (
	"testing"
	"time"
)

func TestDefaultRetryConfig_Values(t *testing.T) {
	rc := DefaultRetryConfig()
	if rc.MaxAttempts != 1 {
		t.Errorf("expected MaxAttempts=1, got %d", rc.MaxAttempts)
	}
	if rc.Delay != 0 {
		t.Errorf("expected Delay=0, got %v", rc.Delay)
	}
}

func TestValidateRetry_Valid(t *testing.T) {
	cases := []RetryConfig{
		{MaxAttempts: 1, Delay: 0},
		{MaxAttempts: 3, Delay: 2 * time.Second},
		{MaxAttempts: 10, Delay: 5 * time.Minute},
	}
	for _, c := range cases {
		if err := ValidateRetry(c); err != nil {
			t.Errorf("unexpected error for %+v: %v", c, err)
		}
	}
}

func TestValidateRetry_ZeroAttempts(t *testing.T) {
	err := ValidateRetry(RetryConfig{MaxAttempts: 0, Delay: 0})
	if err == nil {
		t.Fatal("expected error for MaxAttempts=0")
	}
}

func TestValidateRetry_TooManyAttempts(t *testing.T) {
	err := ValidateRetry(RetryConfig{MaxAttempts: 11, Delay: 0})
	if err == nil {
		t.Fatal("expected error for MaxAttempts=11")
	}
}

func TestValidateRetry_NegativeDelay(t *testing.T) {
	err := ValidateRetry(RetryConfig{MaxAttempts: 1, Delay: -time.Second})
	if err == nil {
		t.Fatal("expected error for negative delay")
	}
}

func TestValidateRetry_DelayTooLarge(t *testing.T) {
	err := ValidateRetry(RetryConfig{MaxAttempts: 1, Delay: 6 * time.Minute})
	if err == nil {
		t.Fatal("expected error for delay > 5m")
	}
}
