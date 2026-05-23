package executor

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

// TimeoutRunner wraps command execution with a configurable timeout.
type TimeoutRunner struct {
	defaultTimeout time.Duration
}

// NewTimeoutRunner creates a TimeoutRunner with the given default timeout.
// A zero or negative duration means no timeout is applied.
func NewTimeoutRunner(timeout time.Duration) *TimeoutRunner {
	return &TimeoutRunner{defaultTimeout: timeout}
}

// RunWithTimeout executes the given shell command string under a timeout context.
// If timeout is zero the runner's default is used; if both are zero the command
// runs without a deadline.
func (t *TimeoutRunner) RunWithTimeout(shell, command string, timeout time.Duration) ([]byte, error) {
	effective := timeout
	if effective <= 0 {
		effective = t.defaultTimeout
	}

	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	if effective > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), effective)
		defer cancel()
	} else {
		ctx = context.Background()
	}

	cmd := exec.CommandContext(ctx, shell, "-c", command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return out, fmt.Errorf("command timed out after %s: %w", effective, ctx.Err())
		}
		return out, fmt.Errorf("command failed: %w", err)
	}
	return out, nil
}
