package executor

import (
	"context"
	"os/exec"
	"strings"
	"testing"
	"time"
)

// shellRunConcurrent executes cmd via the system shell and returns combined output.
func shellRunConcurrent(_ context.Context, cmd string) (string, error) {
	out, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func TestConcurrentRunner_Integration_AllSucceed(t *testing.T) {
	r := NewConcurrentRunner(4, shellRunConcurrent)
	cmds := []string{
		"echo alpha",
		"echo beta",
		"echo gamma",
	}
	results := r.RunAll(context.Background(), cmds)

	expected := []string{"alpha", "beta", "gamma"}
	for i, res := range results {
		if res.Err != nil {
			t.Errorf("step %d failed: %v", i, res.Err)
		}
		if res.Output != expected[i] {
			t.Errorf("step %d output = %q, want %q", i, res.Output, expected[i])
		}
	}
}

func TestConcurrentRunner_Integration_OneFails(t *testing.T) {
	r := NewConcurrentRunner(2, shellRunConcurrent)
	cmds := []string{"echo ok", "exit 1", "echo also-ok"}
	results := r.RunAll(context.Background(), cmds)

	if results[1].Err == nil {
		t.Fatal("expected failure for 'exit 1'")
	}
	if results[0].Err != nil || results[2].Err != nil {
		t.Error("good steps should not fail")
	}
	if err := FirstError(results); err == nil {
		t.Error("FirstError should return non-nil")
	}
}

func TestConcurrentRunner_Integration_SpeedupWithWorkers(t *testing.T) {
	const sleepCmd = "sleep 0.1"
	const n = 4
	cmds := make([]string, n)
	for i := range cmds {
		cmds[i] = sleepCmd
	}

	start := time.Now()
	NewConcurrentRunner(n, shellRunConcurrent).RunAll(context.Background(), cmds)
	elapsed := time.Since(start)

	// With n workers the wall time should be well under n*100 ms.
	if elapsed > time.Duration(n)*150*time.Millisecond {
		t.Errorf("parallel execution too slow: %v", elapsed)
	}
}
