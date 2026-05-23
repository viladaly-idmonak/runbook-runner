package config

import "testing"

func TestDefaultConcurrencyConfig_Values(t *testing.T) {
	cfg := DefaultConcurrencyConfig()
	if cfg.Workers != DefaultWorkers {
		t.Errorf("Workers = %d, want %d", cfg.Workers, DefaultWorkers)
	}
	if !cfg.FailFast {
		t.Error("FailFast should default to true")
	}
}

func TestValidateConcurrency_Valid(t *testing.T) {
	cases := []ConcurrencyConfig{
		{Workers: 1, FailFast: false},
		{Workers: 8, FailFast: true},
		{Workers: MaxWorkers, FailFast: false},
	}
	for _, c := range cases {
		if err := ValidateConcurrency(c); err != nil {
			t.Errorf("unexpected error for %+v: %v", c, err)
		}
	}
}

func TestValidateConcurrency_ZeroWorkers(t *testing.T) {
	err := ValidateConcurrency(ConcurrencyConfig{Workers: 0})
	if err == nil {
		t.Fatal("expected error for Workers=0")
	}
}

func TestValidateConcurrency_NegativeWorkers(t *testing.T) {
	err := ValidateConcurrency(ConcurrencyConfig{Workers: -1})
	if err == nil {
		t.Fatal("expected error for Workers=-1")
	}
}

func TestValidateConcurrency_TooManyWorkers(t *testing.T) {
	err := ValidateConcurrency(ConcurrencyConfig{Workers: MaxWorkers + 1})
	if err == nil {
		t.Fatalf("expected error for Workers=%d", MaxWorkers+1)
	}
}
