package config

import (
	"testing"
	"time"
)

func TestDefaultRateConfig_Values(t *testing.T) {
	cfg := DefaultRateConfig()
	if cfg.Enabled {
		t.Error("expected Enabled to be false")
	}
	if cfg.MaxPerMinute != 60 {
		t.Errorf("expected MaxPerMinute 60, got %d", cfg.MaxPerMinute)
	}
	if cfg.Burst != 5 {
		t.Errorf("expected Burst 5, got %d", cfg.Burst)
	}
	if cfg.MinStepDelay != 0 {
		t.Errorf("expected MinStepDelay 0, got %v", cfg.MinStepDelay)
	}
}

func TestValidateRate_DisabledSkipsValidation(t *testing.T) {
	cfg := RateConfig{Enabled: false, MaxPerMinute: -1}
	if err := ValidateRate(cfg); err != nil {
		t.Errorf("expected no error for disabled config, got %v", err)
	}
}

func TestValidateRate_Valid(t *testing.T) {
	cfg := RateConfig{Enabled: true, MaxPerMinute: 30, Burst: 3, MinStepDelay: 500 * time.Millisecond}
	if err := ValidateRate(cfg); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateRate_ZeroMaxPerMinute(t *testing.T) {
	cfg := RateConfig{Enabled: true, MaxPerMinute: 0, Burst: 1}
	if err := ValidateRate(cfg); err == nil {
		t.Error("expected error for zero max_per_minute")
	}
}

func TestValidateRate_ExceedsMaxPerMinute(t *testing.T) {
	cfg := RateConfig{Enabled: true, MaxPerMinute: 9999, Burst: 1}
	if err := ValidateRate(cfg); err == nil {
		t.Error("expected error for max_per_minute > 3600")
	}
}

func TestValidateRate_BurstExceedsMaxPerMinute(t *testing.T) {
	cfg := RateConfig{Enabled: true, MaxPerMinute: 10, Burst: 20}
	if err := ValidateRate(cfg); err == nil {
		t.Error("expected error when burst exceeds max_per_minute")
	}
}

func TestValidateRate_NegativeMinStepDelay(t *testing.T) {
	cfg := RateConfig{Enabled: true, MaxPerMinute: 10, Burst: 2, MinStepDelay: -1 * time.Second}
	if err := ValidateRate(cfg); err == nil {
		t.Error("expected error for negative min_step_delay")
	}
}
