package executor

import (
	"errors"

	"github.com/your-org/runbook-runner/internal/config"
	"github.com/your-org/runbook-runner/internal/parser"
)

// ErrStepLimitExceeded is returned when the runbook exceeds the configured
// maximum number of steps and FailOnExceed is true.
var ErrStepLimitExceeded = errors.New("step limit exceeded")

// StepLimiter trims or rejects a step list based on StepLimitConfig.
type StepLimiter struct {
	cfg config.StepLimitConfig
}

// NewStepLimiter creates a StepLimiter from the provided configuration.
func NewStepLimiter(cfg config.StepLimitConfig) *StepLimiter {
	return &StepLimiter{cfg: cfg}
}

// Apply returns the (possibly truncated) slice of steps together with an
// optional error. When the limit is disabled every step is returned unchanged.
// When FailOnExceed is true and the step count exceeds MaxSteps an error is
// returned. Otherwise the slice is silently truncated to MaxSteps.
func (l *StepLimiter) Apply(steps []parser.Step) ([]parser.Step, error) {
	if !l.cfg.Enabled {
		return steps, nil
	}
	if len(steps) <= l.cfg.MaxSteps {
		return steps, nil
	}
	if l.cfg.FailOnExceed {
		return nil, ErrStepLimitExceeded
	}
	return steps[:l.cfg.MaxSteps], nil
}

// Remaining returns how many additional steps could be run given the current
// count. Returns -1 when the limiter is disabled (unlimited).
func (l *StepLimiter) Remaining(executed int) int {
	if !l.cfg.Enabled {
		return -1
	}
	rem := l.cfg.MaxSteps - executed
	if rem < 0 {
		return 0
	}
	return rem
}
