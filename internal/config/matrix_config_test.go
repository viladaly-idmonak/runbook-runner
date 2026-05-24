package config

import "testing"

func TestDefaultMatrixConfig_Values(t *testing.T) {
	cfg := DefaultMatrixConfig()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if len(cfg.Vars) != 0 {
		t.Errorf("expected empty Vars, got %v", cfg.Vars)
	}
	if cfg.MaxExpansions != 50 {
		t.Errorf("expected MaxExpansions=50, got %d", cfg.MaxExpansions)
	}
}

func TestValidateMatrix_DisabledSkipsValidation(t *testing.T) {
	cfg := MatrixConfig{Enabled: false} // no vars — still valid when disabled
	if err := ValidateMatrix(cfg); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateMatrix_EnabledNoVars(t *testing.T) {
	cfg := MatrixConfig{Enabled: true, Vars: map[string][]string{}, MaxExpansions: 50}
	if err := ValidateMatrix(cfg); err == nil {
		t.Error("expected error for enabled matrix with no vars")
	}
}

func TestValidateMatrix_ValidConfig(t *testing.T) {
	cfg := MatrixConfig{
		Enabled:       true,
		Vars:          map[string][]string{"ENV": {"staging", "prod"}},
		MaxExpansions: 50,
	}
	if err := ValidateMatrix(cfg); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateMatrix_EmptyVarKey(t *testing.T) {
	cfg := MatrixConfig{
		Enabled:       true,
		Vars:          map[string][]string{"": {"a"}},
		MaxExpansions: 50,
	}
	if err := ValidateMatrix(cfg); err == nil {
		t.Error("expected error for empty var key")
	}
}

func TestValidateMatrix_EmptyVarValues(t *testing.T) {
	cfg := MatrixConfig{
		Enabled:       true,
		Vars:          map[string][]string{"ENV": {}},
		MaxExpansions: 50,
	}
	if err := ValidateMatrix(cfg); err == nil {
		t.Error("expected error for var with no values")
	}
}

func TestValidateMatrix_EmptyValueInList(t *testing.T) {
	cfg := MatrixConfig{
		Enabled:       true,
		Vars:          map[string][]string{"ENV": {"staging", ""}},
		MaxExpansions: 50,
	}
	if err := ValidateMatrix(cfg); err == nil {
		t.Error("expected error for empty string in var values")
	}
}

func TestValidateMatrix_NegativeMaxExpansions(t *testing.T) {
	cfg := MatrixConfig{
		Enabled:       true,
		Vars:          map[string][]string{"ENV": {"prod"}},
		MaxExpansions: -1,
	}
	if err := ValidateMatrix(cfg); err == nil {
		t.Error("expected error for negative max_expansions")
	}
}

func TestValidateMatrix_ExceedsMaxExpansions(t *testing.T) {
	cfg := MatrixConfig{
		Enabled: true,
		Vars: map[string][]string{
			"ENV":    {"staging", "prod"},
			"REGION": {"us", "eu", "ap"},
		},
		MaxExpansions: 3, // 2*3=6 > 3
	}
	if err := ValidateMatrix(cfg); err == nil {
		t.Error("expected error when expansion count exceeds max")
	}
}

func TestValidateMatrix_ZeroMaxExpansionsMeansNoLimit(t *testing.T) {
	cfg := MatrixConfig{
		Enabled: true,
		Vars: map[string][]string{
			"ENV":    {"staging", "prod"},
			"REGION": {"us", "eu", "ap"},
		},
		MaxExpansions: 0,
	}
	if err := ValidateMatrix(cfg); err != nil {
		t.Errorf("unexpected error with unlimited expansions: %v", err)
	}
}
