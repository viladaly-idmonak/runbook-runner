package config_test

import (
	"testing"

	"github.com/your-org/runbook-runner/internal/config"
)

func TestDefaultDryRunConfig_Values(t *testing.T) {
	cfg := config.DefaultDryRunConfig()

	if cfg.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if !cfg.PrintCommands {
		t.Error("expected PrintCommands to be true by default")
	}
	if cfg.ExitOnFirstSkip {
		t.Error("expected ExitOnFirstSkip to be false by default")
	}
}

func TestValidateDryRun_DisabledNoExitOnFirstSkip(t *testing.T) {
	cfg := config.DryRunConfig{
		Enabled:         false,
		PrintCommands:   true,
		ExitOnFirstSkip: false,
	}
	if err := config.ValidateDryRun(cfg); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidateDryRun_EnabledWithExitOnFirstSkip(t *testing.T) {
	cfg := config.DryRunConfig{
		Enabled:         true,
		PrintCommands:   true,
		ExitOnFirstSkip: true,
	}
	if err := config.ValidateDryRun(cfg); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidateDryRun_DisabledButExitOnFirstSkip(t *testing.T) {
	cfg := config.DryRunConfig{
		Enabled:         false,
		PrintCommands:   false,
		ExitOnFirstSkip: true,
	}
	err := config.ValidateDryRun(cfg)
	if err == nil {
		t.Fatal("expected error when exit_on_first_skip is set without enabled")
	}
	expected := "dry_run: exit_on_first_skip requires enabled to be true"
	if err.Error() != expected {
		t.Errorf("unexpected error message: %q", err.Error())
	}
}

func TestValidateDryRun_EnabledPrintCommandsFalse(t *testing.T) {
	cfg := config.DryRunConfig{
		Enabled:         true,
		PrintCommands:   false,
		ExitOnFirstSkip: false,
	}
	if err := config.ValidateDryRun(cfg); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}
