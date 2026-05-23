package executor

import (
	"strings"
	"testing"
	"time"
)

// realRunner executes commands via the shell for integration tests.
type realRunner struct{}

func (r *realRunner) Run(cmd string) (string, error) {
	return shellRun(cmd)
}

// shellRun is a thin wrapper so the integration test doesn't depend on runner internals.
func shellRun(cmd string) (string, error) {
	opts := DefaultOptions()
	opts.DryRun = false
	runner := New(opts)
	return runner.Run(cmd)
}

func TestRetryRunner_Integration_EventualSuccess(t *testing.T) {
	// Use a counter file to fail the first two attempts and succeed on the third.
	countFile := t.TempDir() + "/count"
	cmd := `
		COUNT=$(cat "` + countFile + `" 2>/dev/null || echo 0)
		COUNT=$((COUNT+1))
		echo $COUNT > "` + countFile + `"
		if [ "$COUNT" -lt 3 ]; then exit 1; fi
		echo "success"
	`
	rr := NewRetryRunner(&realRunner{}, RetryPolicy{MaxAttempts: 5, Delay: 0})
	out, err := rr.Run(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "success") {
		t.Errorf("expected 'success' in output, got %q", out)
	}
}

func TestRetryRunner_Integration_AlwaysFails(t *testing.T) {
	rr := NewRetryRunner(&realRunner{}, RetryPolicy{MaxAttempts: 3, Delay: 0})
	_, err := rr.Run("exit 1")
	if err == nil {
		t.Fatal("expected error when command always fails")
	}
	if !strings.Contains(err.Error(), "3 attempt(s)") {
		t.Errorf("error message should mention attempt count: %v", err)
	}
}

func TestRetryRunner_Integration_DelayIsRespected(t *testing.T) {
	rr := NewRetryRunner(&realRunner{}, RetryPolicy{MaxAttempts: 2, Delay: 50 * time.Millisecond})
	start := time.Now()
	rr.Run("exit 1") //nolint:errcheck
	elapsed := time.Since(start)
	if elapsed < 50*time.Millisecond {
		t.Errorf("expected at least 50ms delay, got %v", elapsed)
	}
}
