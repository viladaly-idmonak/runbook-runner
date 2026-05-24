package executor

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/your-org/runbook-runner/internal/config"
)

func TestAuditLogger_Disabled_NoFile(t *testing.T) {
	cfg := config.DefaultAuditConfig()
	cfg.Enabled = false

	logger, err := NewAuditLogger(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer logger.Close()

	if logger.file != nil {
		t.Error("expected file to be nil when auditing is disabled")
	}
}

func TestAuditLogger_Disabled_RecordIsNoop(t *testing.T) {
	cfg := config.DefaultAuditConfig()
	cfg.Enabled = false

	logger, _ := NewAuditLogger(cfg)
	defer logger.Close()

	if err := logger.Record(AuditEntry{StepName: "step1", Status: "ok"}); err != nil {
		t.Fatalf("expected no error from disabled logger, got: %v", err)
	}
}

func TestAuditLogger_TextFormat_WritesLine(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	cfg := config.DefaultAuditConfig()
	cfg.Enabled = true
	cfg.FilePath = path
	cfg.Format = "text"

	logger, err := NewAuditLogger(cfg)
	if err != nil {
		t.Fatalf("NewAuditLogger: %v", err)
	}

	_ = logger.Record(AuditEntry{
		RunbookTitle: "deploy",
		StepName:     "run migrations",
		Status:       "ok",
		Duration:     "1.2s",
	})
	logger.Close()

	data, _ := os.ReadFile(path)
	line := string(data)
	if !strings.Contains(line, "step=\"run migrations\"") {
		t.Errorf("expected step name in log line, got: %s", line)
	}
	if !strings.Contains(line, "status=ok") {
		t.Errorf("expected status in log line, got: %s", line)
	}
}

func TestAuditLogger_JSONFormat_ValidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.json")

	cfg := config.DefaultAuditConfig()
	cfg.Enabled = true
	cfg.FilePath = path
	cfg.Format = "json"

	logger, err := NewAuditLogger(cfg)
	if err != nil {
		t.Fatalf("NewAuditLogger: %v", err)
	}

	_ = logger.Record(AuditEntry{
		RunbookTitle: "deploy",
		StepName:     "health check",
		Status:       "fail",
		Duration:     "0.3s",
		Error:        "exit status 1",
	})
	logger.Close()

	f, _ := os.Open(path)
	defer f.Close()

	var entry AuditEntry
	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		t.Fatal("expected at least one line in audit log")
	}
	if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if entry.StepName != "health check" {
		t.Errorf("expected step_name=\"health check\", got %q", entry.StepName)
	}
	if entry.Status != "fail" {
		t.Errorf("expected status=\"fail\", got %q", entry.Status)
	}
}
