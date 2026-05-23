package executor

import (
	"runtime"
	"testing"

	"github.com/runbook-runner/internal/parser"
)

func sampleRunbook(cmds []string) *parser.Runbook {
	rb := &parser.Runbook{Title: "Test Runbook"}
	for i, c := range cmds {
		rb.Steps = append(rb.Steps, parser.Step{
			Name:    fmt.Sprintf("step-%d", i+1),
			Command: c,
		})
	}
	return rb
}

func TestRunner_DryRun(t *testing.T) {
	opts := DefaultOptions()
	opts.DryRun = true
	r := New(opts)
	rb := &parser.Runbook{
		Title: "dry run book",
		Steps: []parser.Step{
			{Name: "echo step", Command: "echo hello"},
		},
	}
	if err := r.Run(rb); err != nil {
		t.Fatalf("dry run should not fail: %v", err)
	}
	if got := r.Results()[0].Output; got != "(dry-run)" {
		t.Errorf("expected (dry-run), got %q", got)
	}
}

func TestRunner_SuccessfulSteps(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping shell test on windows")
	}
	r := New(DefaultOptions())
	rb := &parser.Runbook{
		Title: "success book",
		Steps: []parser.Step{
			{Name: "true step", Command: "true"},
			{Name: "echo step", Command: "echo ok"},
		},
	}
	if err := r.Run(rb); err != nil {
		t.Fatalf("expected success, got: %v", err)
	}
	if len(r.Results()) != 2 {
		t.Errorf("expected 2 results, got %d", len(r.Results()))
	}
}

func TestRunner_FailingStep(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping shell test on windows")
	}
	r := New(DefaultOptions())
	rb := &parser.Runbook{
		Title: "fail book",
		Steps: []parser.Step{
			{Name: "good step", Command: "true"},
			{Name: "bad step", Command: "false"},
			{Name: "unreachable", Command: "echo should not run"},
		},
	}
	if err := r.Run(rb); err == nil {
		t.Fatal("expected error from failing step")
	}
	if len(r.Results()) != 2 {
		t.Errorf("expected 2 results (stopped at failure), got %d", len(r.Results()))
	}
}

func TestRunner_EmptyRunbook(t *testing.T) {
	r := New(DefaultOptions())
	rb := &parser.Runbook{Title: "empty"}
	if err := r.Run(rb); err != nil {
		t.Fatalf("empty runbook should succeed: %v", err)
	}
}
