package executor

import (
	"fmt"
	"time"
)

// RetryPolicy defines how a step should be retried on failure.
type RetryPolicy struct {
	MaxAttempts int
	Delay       time.Duration
}

// DefaultRetryPolicy returns a policy with no retries.
func DefaultRetryPolicy() RetryPolicy {
	return RetryPolicy{
		MaxAttempts: 1,
		Delay:       0,
	}
}

// RetryRunner wraps a CommandRunner and retries failed commands.
type RetryRunner struct {
	inner  CommandRunner
	policy RetryPolicy
}

// NewRetryRunner creates a RetryRunner with the given policy.
func NewRetryRunner(inner CommandRunner, policy RetryPolicy) *RetryRunner {
	if policy.MaxAttempts < 1 {
		policy.MaxAttempts = 1
	}
	return &RetryRunner{inner: inner, policy: policy}
}

// Run executes the command, retrying up to MaxAttempts times on failure.
func (r *RetryRunner) Run(cmd string) (string, error) {
	var (
		out string
		err error
	)
	for attempt := 1; attempt <= r.policy.MaxAttempts; attempt++ {
		out, err = r.inner.Run(cmd)
		if err == nil {
			return out, nil
		}
		if attempt < r.policy.MaxAttempts && r.policy.Delay > 0 {
			time.Sleep(r.policy.Delay)
		}
	}
	return out, fmt.Errorf("command failed after %d attempt(s): %w", r.policy.MaxAttempts, err)
}
