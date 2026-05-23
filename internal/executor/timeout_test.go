package executor

import (
	"strings"
	"testing"
	"time"
)

const testShell = "/bin/sh"

func TestTimeoutRunner_SuccessNoTimeout(t *testing.T) {
	tr := NewTimeoutRunner(0)
	out, err := tr.RunWithTimeout(testShell, "echo hello", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(out), "hello") {
		t.Errorf("expected output to contain 'hello', got %q", string(out))
	}
}

func TestTimeoutRunner_SuccessWithDefaultTimeout(t *testing.T) {
	tr := NewTimeoutRunner(5 * time.Second)
	out, err := tr.RunWithTimeout(testShell, "echo world", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(out), "world") {
		t.Errorf("expected 'world' in output, got %q", string(out))
	}
}

func TestTimeoutRunner_PerCallTimeoutOverridesDefault(t *testing.T) {
	tr := NewTimeoutRunner(30 * time.Second)
	// per-call timeout of 100ms should expire before the sleep finishes
	_, err := tr.RunWithTimeout(testShell, "sleep 5", 100*time.Millisecond)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if !strings.Contains(err.Error(), "timed out") {
		t.Errorf("expected 'timed out' in error, got: %v", err)
	}
}

func TestTimeoutRunner_DefaultTimeoutExpires(t *testing.T) {
	tr := NewTimeoutRunner(100 * time.Millisecond)
	_, err := tr.RunWithTimeout(testShell, "sleep 5", 0)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if !strings.Contains(err.Error(), "timed out") {
		t.Errorf("expected 'timed out' in error, got: %v", err)
	}
}

func TestTimeoutRunner_CommandFailureNotTimeout(t *testing.T) {
	tr := NewTimeoutRunner(5 * time.Second)
	_, err := tr.RunWithTimeout(testShell, "exit 1", 0)
	if err == nil {
		t.Fatal("expected error for failing command")
	}
	if strings.Contains(err.Error(), "timed out") {
		t.Errorf("error should not indicate timeout, got: %v", err)
	}
}
