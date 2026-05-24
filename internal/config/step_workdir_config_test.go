package config_test

import (
	"testing"

	"github.com/your-org/runbook-runner/internal/config"
)

func TestDefaultStepWorkdirConfig_Values(t *testing.T) {
	c := config.DefaultStepWorkdirConfig()
	if c.Enabled {
		t.Error("expected Enabled=false")
	}
	if c.Default != "" {
		t.Errorf("expected empty Default, got %q", c.Default)
	}
	if len(c.Overrides) != 0 {
		t.Errorf("expected empty Overrides, got %v", c.Overrides)
	}
}

func TestValidateStepWorkdir_DisabledSkipsValidation(t *testing.T) {
	c := config.StepWorkdirConfig{Enabled: false, Default: "relative/path"}
	if err := config.ValidateStepWorkdir(c); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateStepWorkdir_ValidConfig(t *testing.T) {
	c := config.StepWorkdirConfig{
		Enabled:   true,
		Default:   "/tmp/runbook",
		Overrides: map[string]string{"deploy": "/var/deploy"},
	}
	if err := config.ValidateStepWorkdir(c); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateStepWorkdir_RelativeDefault(t *testing.T) {
	c := config.StepWorkdirConfig{Enabled: true, Default: "relative/dir"}
	if err := config.ValidateStepWorkdir(c); err == nil {
		t.Error("expected error for relative default path")
	}
}

func TestValidateStepWorkdir_EmptyStepName(t *testing.T) {
	c := config.StepWorkdirConfig{
		Enabled:   true,
		Overrides: map[string]string{"": "/tmp/x"},
	}
	if err := config.ValidateStepWorkdir(c); err == nil {
		t.Error("expected error for empty step name in overrides")
	}
}

func TestValidateStepWorkdir_RelativeOverride(t *testing.T) {
	c := config.StepWorkdirConfig{
		Enabled:   true,
		Overrides: map[string]string{"build": "relative/path"},
	}
	if err := config.ValidateStepWorkdir(c); err == nil {
		t.Error("expected error for relative override path")
	}
}

func TestDirForStep_Disabled(t *testing.T) {
	c := config.StepWorkdirConfig{Enabled: false, Default: "/tmp"}
	if got := c.DirForStep("any"); got != "" {
		t.Errorf("expected empty string when disabled, got %q", got)
	}
}

func TestDirForStep_UsesOverride(t *testing.T) {
	c := config.StepWorkdirConfig{
		Enabled:   true,
		Default:   "/default",
		Overrides: map[string]string{"deploy": "/deploy"},
	}
	if got := c.DirForStep("deploy"); got != "/deploy" {
		t.Errorf("expected /deploy, got %q", got)
	}
}

func TestDirForStep_FallsBackToDefault(t *testing.T) {
	c := config.StepWorkdirConfig{
		Enabled:   true,
		Default:   "/default",
		Overrides: map[string]string{},
	}
	if got := c.DirForStep("unknown"); got != "/default" {
		t.Errorf("expected /default, got %q", got)
	}
}
