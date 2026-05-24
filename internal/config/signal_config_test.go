package config_test

import (
	"testing"

	"github.com/your-org/runbook-runner/internal/config"
)

func TestDefaultSignalConfig_Values(t *testing.T) {
	c := config.DefaultSignalConfig()
	if !c.GracefulShutdown {
		t.Error("expected GracefulShutdown to be true by default")
	}
	if c.RunRollbackOnSignal {
		t.Error("expected RunRollbackOnSignal to be false by default")
	}
	if c.NotifySignal != "SIGINT" {
		t.Errorf("expected NotifySignal=SIGINT, got %q", c.NotifySignal)
	}
}

func TestValidateSignal_ValidSIGINT(t *testing.T) {
	c := config.DefaultSignalConfig()
	if err := config.ValidateSignal(c); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateSignal_ValidSIGTERM(t *testing.T) {
	c := config.DefaultSignalConfig()
	c.NotifySignal = "SIGTERM"
	if err := config.ValidateSignal(c); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateSignal_EmptyNotifySignal(t *testing.T) {
	c := config.DefaultSignalConfig()
	c.NotifySignal = ""
	if err := config.ValidateSignal(c); err == nil {
		t.Error("expected error for empty notify_signal")
	}
}

func TestValidateSignal_UnknownSignal(t *testing.T) {
	c := config.DefaultSignalConfig()
	c.NotifySignal = "SIGUSR1"
	if err := config.ValidateSignal(c); err == nil {
		t.Error("expected error for unknown signal")
	}
}

func TestValidateSignal_RollbackOnSignalEnabled(t *testing.T) {
	c := config.DefaultSignalConfig()
	c.RunRollbackOnSignal = true
	if err := config.ValidateSignal(c); err != nil {
		t.Errorf("unexpected error with rollback enabled: %v", err)
	}
}
