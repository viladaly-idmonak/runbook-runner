package executor_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/runbook-runner/internal/executor"
)

// TestStepLogger_MultipleEntries verifies sequential logging produces one line per step.
func TestStepLogger_MultipleEntries(t *testing.T) {
	var buf bytes.Buffer
	sl := executor.NewStepLogger(&buf, false)

	steps := []executor.StepLog{
		{
			StepName:  "Step 1",
			Command:   "echo hello",
			StartedAt: time.Now(),
			Duration:  10 * time.Millisecond,
		},
		{
			StepName:  "Step 2",
			Command:   "echo world",
			StartedAt: time.Now(),
			Duration:  20 * time.Millisecond,
		},
		{
			StepName:  "Step 3",
			Command:   "exit 1",
			StartedAt: time.Now(),
			Duration:  5 * time.Millisecond,
			Skipped:   true,
		},
	}

	for _, s := range steps {
		sl.Log(s)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != len(steps) {
		t.Errorf("expected %d log lines, got %d:\n%s", len(steps), len(lines), buf.String())
	}

	if !strings.Contains(lines[0], "Step 1") {
		t.Errorf("line 0 should mention Step 1, got: %s", lines[0])
	}
	if !strings.Contains(lines[2], "[SKIP]") {
		t.Errorf("line 2 should be SKIP, got: %s", lines[2])
	}
}
