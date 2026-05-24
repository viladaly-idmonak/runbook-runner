package executor

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/your-org/runbook-runner/internal/config"
)

// AuditEntry represents a single recorded step execution event.
type AuditEntry struct {
	Timestamp time.Time `json:"timestamp"`
	RunbookTitle string    `json:"runbook_title"`
	StepName     string    `json:"step_name"`
	Command      string    `json:"command"`
	Status       string    `json:"status"`
	Duration     string    `json:"duration"`
	Output       string    `json:"output,omitempty"`
	Error        string    `json:"error,omitempty"`
}

// AuditLogger writes step execution records to a file.
type AuditLogger struct {
	cfg  config.AuditConfig
	file *os.File
}

// NewAuditLogger creates an AuditLogger. If auditing is disabled, a no-op
// logger is returned (file == nil).
func NewAuditLogger(cfg config.AuditConfig) (*AuditLogger, error) {
	if !cfg.Enabled {
		return &AuditLogger{cfg: cfg}, nil
	}

	f, err := os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("audit: open log file %q: %w", cfg.FilePath, err)
	}
	return &AuditLogger{cfg: cfg, file: f}, nil
}

// Record writes an AuditEntry. It is a no-op when auditing is disabled.
func (a *AuditLogger) Record(entry AuditEntry) error {
	if !a.cfg.Enabled || a.file == nil {
		return nil
	}

	entry.Timestamp = time.Now().UTC()

	switch a.cfg.Format {
	case "json":
		data, err := json.Marshal(entry)
		if err != nil {
			return fmt.Errorf("audit: marshal entry: %w", err)
		}
		_, err = fmt.Fprintln(a.file, string(data))
		return err
	default: // "text"
		_, err := fmt.Fprintf(a.file, "[%s] runbook=%q step=%q status=%s duration=%s\n",
			entry.Timestamp.Format(time.RFC3339),
			entry.RunbookTitle,
			entry.StepName,
			entry.Status,
			entry.Duration,
		)
		return err
	}
}

// Close releases the underlying file handle.
func (a *AuditLogger) Close() error {
	if a.file != nil {
		return a.file.Close()
	}
	return nil
}
