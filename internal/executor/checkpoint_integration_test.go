package executor_test

import (
	"testing"

	"github.com/your-org/runbook-runner/internal/executor"
)

// TestCheckpointStore_Integration_ResumeSkipsCompleted verifies that a second
// pass over steps correctly skips those already recorded in the checkpoint.
func TestCheckpointStore_Integration_ResumeSkipsCompleted(t *testing.T) {
	dir := t.TempDir()
	store := executor.NewCheckpointStore(dir, "integration-runbook")

	steps := []string{"install-deps", "build", "test", "deploy"}

	// Simulate first run: complete first two steps then "fail".
	for _, s := range steps[:2] {
		if err := store.MarkDone(s); err != nil {
			t.Fatalf("MarkDone %q: %v", s, err)
		}
	}

	// Second run: collect which steps would be skipped vs executed.
	var skipped, executed []string
	for _, s := range steps {
		if store.IsDone(s) {
			skipped = append(skipped, s)
		} else {
			executed = append(executed, s)
		}
	}

	if len(skipped) != 2 {
		t.Errorf("expected 2 skipped steps, got %d: %v", len(skipped), skipped)
	}
	if len(executed) != 2 {
		t.Errorf("expected 2 executed steps, got %d: %v", len(executed), executed)
	}
	if executed[0] != "test" || executed[1] != "deploy" {
		t.Errorf("unexpected executed steps: %v", executed)
	}
}

// TestCheckpointStore_Integration_FullResetAllowsRerun verifies Reset clears state.
func TestCheckpointStore_Integration_FullResetAllowsRerun(t *testing.T) {
	dir := t.TempDir()
	store := executor.NewCheckpointStore(dir, "reset-runbook")

	steps := []string{"step-a", "step-b", "step-c"}
	for _, s := range steps {
		_ = store.MarkDone(s)
	}
	if err := store.Reset(); err != nil {
		t.Fatalf("Reset: %v", err)
	}
	for _, s := range steps {
		if store.IsDone(s) {
			t.Errorf("step %q should not be done after reset", s)
		}
	}
}
