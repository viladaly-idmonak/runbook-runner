package config_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/yourorg/runbook-runner/internal/config"
)

func TestDefault_Values(t *testing.T) {
	cfg := config.Default()
	if cfg.DryRun {
		t.Error("expected DryRun to be false")
	}
	if !cfg.StopOnFailure {
		t.Error("expected StopOnFailure to be true")
	}
	if cfg.ShellPath != "/bin/sh" {
		t.Errorf("unexpected ShellPath: %q", cfg.ShellPath)
	}
	if cfg.OutputFormat != "text" {
		t.Errorf("unexpected OutputFormat: %q", cfg.OutputFormat)
	}
	if cfg.StepTimeout != 30*time.Second {
		t.Errorf("unexpected StepTimeout: %v", cfg.StepTimeout)
	}
}

func TestLoadFile_Valid(t *testing.T) {
	data := map[string]interface{}{
		"dry_run":       true,
		"output_format": "json",
		"step_timeout":  int64(10 * time.Second),
		"shell_path":    "/bin/bash",
	}
	f, err := os.CreateTemp("", "runbook-cfg-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	if err := json.NewEncoder(f).Encode(data); err != nil {
		t.Fatal(err)
	}
	f.Close()

	cfg, err := config.LoadFile(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.DryRun {
		t.Error("expected DryRun true")
	}
	if cfg.OutputFormat != "json" {
		t.Errorf("expected output_format json, got %q", cfg.OutputFormat)
	}
	if cfg.ShellPath != "/bin/bash" {
		t.Errorf("expected /bin/bash, got %q", cfg.ShellPath)
	}
}

func TestLoadFile_Missing(t *testing.T) {
	_, err := config.LoadFile("/nonexistent/path/config.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestValidate_InvalidOutputFormat(t *testing.T) {
	cfg := config.Default()
	cfg.OutputFormat = "xml"
	if err := cfg.Validate(); err != config.ErrInvalidOutputFormat {
		t.Errorf("expected ErrInvalidOutputFormat, got %v", err)
	}
}

func TestValidate_EmptyShellPath(t *testing.T) {
	cfg := config.Default()
	cfg.ShellPath = ""
	if err := cfg.Validate(); err != config.ErrEmptyShellPath {
		t.Errorf("expected ErrEmptyShellPath, got %v", err)
	}
}

func TestValidate_InvalidTimeout(t *testing.T) {
	cfg := config.Default()
	cfg.StepTimeout = 0
	if err := cfg.Validate(); err != config.ErrInvalidTimeout {
		t.Errorf("expected ErrInvalidTimeout, got %v", err)
	}
}

func TestValidate_Valid(t *testing.T) {
	if err := config.Default().Validate(); err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}
