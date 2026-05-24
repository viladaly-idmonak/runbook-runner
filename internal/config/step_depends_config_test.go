package config_test

import (
	"testing"

	"github.com/your-org/runbook-runner/internal/config"
)

func TestDefaultStepDependsConfig_Values(t *testing.T) {
	c := config.DefaultStepDependsConfig()
	if c.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if !c.FailOnUnmet {
		t.Error("expected FailOnUnmet to be true by default")
	}
	if c.Deps == nil {
		t.Error("expected Deps to be initialised")
	}
}

func TestValidateStepDepends_DisabledSkipsValidation(t *testing.T) {
	c := config.StepDependsConfig{
		Enabled: false,
		Deps:    map[string][]string{"": {"bad"}},
	}
	if err := config.ValidateStepDepends(c); err != nil {
		t.Errorf("expected no error when disabled, got %v", err)
	}
}

func TestValidateStepDepends_Valid(t *testing.T) {
	c := config.StepDependsConfig{
		Enabled: true,
		Deps: map[string][]string{
			"deploy": {"build", "test"},
		},
	}
	if err := config.ValidateStepDepends(c); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateStepDepends_EmptyStepName(t *testing.T) {
	c := config.StepDependsConfig{
		Enabled: true,
		Deps:    map[string][]string{"": {"build"}},
	}
	if err := config.ValidateStepDepends(c); err == nil {
		t.Error("expected error for empty step name")
	}
}

func TestValidateStepDepends_EmptyDepName(t *testing.T) {
	c := config.StepDependsConfig{
		Enabled: true,
		Deps:    map[string][]string{"deploy": {""}},
	}
	if err := config.ValidateStepDepends(c); err == nil {
		t.Error("expected error for empty dependency name")
	}
}

func TestValidateStepDepends_SelfDependency(t *testing.T) {
	c := config.StepDependsConfig{
		Enabled: true,
		Deps:    map[string][]string{"build": {"build"}},
	}
	if err := config.ValidateStepDepends(c); err == nil {
		t.Error("expected error for self-dependency")
	}
}

func TestValidateStepDepends_DuplicateDep(t *testing.T) {
	c := config.StepDependsConfig{
		Enabled: true,
		Deps:    map[string][]string{"deploy": {"build", "build"}},
	}
	if err := config.ValidateStepDepends(c); err == nil {
		t.Error("expected error for duplicate dependency")
	}
}
