package reporter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestJSONReporter_EmptySteps(t *testing.T) {
	var buf bytes.Buffer
	r := NewJSONReporter("empty runbook", &buf)
	if err := r.Flush(); err != nil {
		t.Fatalf("Flush returned error: %v", err)
	}
	var out JSONReport
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if out.Title != "empty runbook" {
		t.Errorf("expected title %q, got %q", "empty runbook", out.Title)
	}
	if len(out.Steps) != 0 {
		t.Errorf("expected 0 steps, got %d", len(out.Steps))
	}
}

func TestJSONReporter_TallyCounts(t *testing.T) {
	var buf bytes.Buffer
	r := NewJSONReporter("tally test", &buf)
	r.Record("step1", "ok", "done", "")
	r.Record("step2", "fail", "", "exit 1")
	r.Record("step3", "skip", "", "")
	r.Record("step4", "ok", "done", "")
	_ = r.Flush()

	var out JSONReport
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out.Passed != 2 {
		t.Errorf("expected 2 passed, got %d", out.Passed)
	}
	if out.Failed != 1 {
		t.Errorf("expected 1 failed, got %d", out.Failed)
	}
	if out.Skipped != 1 {
		t.Errorf("expected 1 skipped, got %d", out.Skipped)
	}
}

func TestJSONReporter_StepFields(t *testing.T) {
	var buf bytes.Buffer
	r := NewJSONReporter("fields test", &buf)
	r.Record("deploy", "fail", "some output", "exit status 2")
	_ = r.Flush()

	var out JSONReport
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	s := out.Steps[0]
	if s.Name != "deploy" {
		t.Errorf("expected name %q, got %q", "deploy", s.Name)
	}
	if s.Status != "fail" {
		t.Errorf("expected status %q, got %q", "fail", s.Status)
	}
	if !strings.Contains(s.Error, "exit status 2") {
		t.Errorf("expected error to contain %q, got %q", "exit status 2", s.Error)
	}
}

func TestJSONReporter_OutputIsIndented(t *testing.T) {
	var buf bytes.Buffer
	r := NewJSONReporter("indent test", &buf)
	_ = r.Flush()
	if !strings.Contains(buf.String(), "\n") {
		t.Error("expected indented (multi-line) JSON output")
	}
}
