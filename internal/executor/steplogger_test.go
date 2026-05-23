package executor

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
)

func baseEntry() StepLog {
	return StepLog{
		StepName:  "Install dependencies",
		Command:   "apt-get install -y curl",
		Output:    "Reading package lists...",
		StartedAt: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		Duration:  350 * time.Millisecond,
	}
}

func TestStepLogger_OKStatus(t *testing.T) {
	var buf bytes.Buffer
	sl := NewStepLogger(&buf, false)
	sl.Log(baseEntry())

	if !strings.Contains(buf.String(), "[OK]") {
		t.Errorf("expected [OK] in output, got: %s", buf.String())
	}
}

func TestStepLogger_FailStatus(t *testing.T) {
	var buf bytes.Buffer
	sl := NewStepLogger(&buf, false)

	entry := baseEntry()
	entry.Err = errors.New("command not found")
	sl.Log(entry)

	out := buf.String()
	if !strings.Contains(out, "[FAIL]") {
		t.Errorf("expected [FAIL] in output, got: %s", out)
	}
	if !strings.Contains(out, "command not found") {
		t.Errorf("expected error message in output, got: %s", out)
	}
}

func TestStepLogger_SkipStatus(t *testing.T) {
	var buf bytes.Buffer
	sl := NewStepLogger(&buf, false)

	entry := baseEntry()
	entry.Skipped = true
	sl.Log(entry)

	if !strings.Contains(buf.String(), "[SKIP]") {
		t.Errorf("expected [SKIP] in output, got: %s", buf.String())
	}
}

func TestStepLogger_VerboseShowsOutput(t *testing.T) {
	var buf bytes.Buffer
	sl := NewStepLogger(&buf, true)
	sl.Log(baseEntry())

	if !strings.Contains(buf.String(), "Reading package lists") {
		t.Errorf("expected command output in verbose mode, got: %s", buf.String())
	}
}

func TestStepLogger_NonVerboseHidesOutput(t *testing.T) {
	var buf bytes.Buffer
	sl := NewStepLogger(&buf, false)
	sl.Log(baseEntry())

	if strings.Contains(buf.String(), "Reading package lists") {
		t.Errorf("expected output hidden in non-verbose mode, got: %s", buf.String())
	}
}

func TestNewStepLogger_NilWriterDefaultsToStdout(t *testing.T) {
	sl := NewStepLogger(nil, false)
	if sl == nil {
		t.Fatal("expected non-nil StepLogger")
	}
	if sl.out == nil {
		t.Error("expected non-nil writer when nil passed")
	}
}
