package config_test

import (
	"testing"

	"github.com/your-org/runbook-runner/internal/config"
)

func TestDefaultCheckpointConfig_Values(t *testing.T) {
	cfg := config.DefaultCheckpointConfig()
	if cfg.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if cfg.Dir != ".runbook-checkpoints" {
		t.Errorf("unexpected default dir: %q", cfg.Dir)
	}
	if !cfg.ResumeOnRestart {
		t.Error("expected ResumeOnRestart to be true by default")
	}
}

func TestValidateCheckpoint_DisabledAlwaysValid(t *testing.T) {
	cfg := config.CheckpointConfig{Enabled: false, Dir: ""}
	if err := config.ValidateCheckpoint(cfg); err != nil {
		t.Errorf("expected no error for disabled checkpoint, got: %v", err)
	}
}

func TestValidateCheckpoint_EnabledWithDir(t *testing.T) {
	cfg := config.CheckpointConfig{Enabled: true, Dir: "/tmp/cp"}
	if err := config.ValidateCheckpoint(cfg); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidateCheckpoint_EnabledEmptyDir(t *testing.T) {
	cfg := config.CheckpointConfig{Enabled: true, Dir: ""}
	if err := config.ValidateCheckpoint(cfg); err == nil {
		t.Error("expected error for empty dir when enabled")
	}
}
