package config_test

import (
	"testing"

	"github.com/your-org/runbook-runner/internal/config"
)

func TestDefaultAuditConfig_Values(t *testing.T) {
	cfg := config.DefaultAuditConfig()

	if cfg.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if cfg.FilePath != "" {
		t.Errorf("expected empty FilePath, got %q", cfg.FilePath)
	}
	if cfg.Format != "text" {
		t.Errorf("expected Format \"text\", got %q", cfg.Format)
	}
	if cfg.IncludeOutput {
		t.Error("expected IncludeOutput to be false by default")
	}
}

func TestValidateAudit_DisabledSkipsValidation(t *testing.T) {
	cfg := config.AuditConfig{
		Enabled:  false,
		FilePath: "",
		Format:   "",
	}
	if err := config.ValidateAudit(cfg); err != nil {
		t.Errorf("expected no error when disabled, got: %v", err)
	}
}

func TestValidateAudit_EnabledWithValidConfig(t *testing.T) {
	cfg := config.AuditConfig{
		Enabled:  true,
		FilePath: "/var/log/runbook-audit.log",
		Format:   "json",
	}
	if err := config.ValidateAudit(cfg); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestValidateAudit_EnabledMissingFilePath(t *testing.T) {
	cfg := config.AuditConfig{
		Enabled:  true,
		FilePath: "",
		Format:   "text",
	}
	if err := config.ValidateAudit(cfg); err == nil {
		t.Error("expected error for missing file_path")
	}
}

func TestValidateAudit_EnabledUnknownFormat(t *testing.T) {
	cfg := config.AuditConfig{
		Enabled:  true,
		FilePath: "/tmp/audit.log",
		Format:   "csv",
	}
	if err := config.ValidateAudit(cfg); err == nil {
		t.Error("expected error for unknown format")
	}
}

func TestValidateAudit_EnabledEmptyFormat(t *testing.T) {
	cfg := config.AuditConfig{
		Enabled:  true,
		FilePath: "/tmp/audit.log",
		Format:   "",
	}
	if err := config.ValidateAudit(cfg); err == nil {
		t.Error("expected error for empty format")
	}
}

func TestValidateAudit_TextFormatIsValid(t *testing.T) {
	cfg := config.AuditConfig{
		Enabled:  true,
		FilePath: "/tmp/audit.log",
		Format:   "text",
	}
	if err := config.ValidateAudit(cfg); err != nil {
		t.Errorf("expected no error for text format, got: %v", err)
	}
}
