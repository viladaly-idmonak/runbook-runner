package executor

import (
	"time"

	"github.com/your-org/runbook-runner/internal/config"
)

// StepRetryResolver derives the effective RetryPolicy for a named step by
// merging the global RetryConfig with any per-step StepRetryConfig overrides.
type StepRetryResolver struct {
	global   config.RetryConfig
	perStep  config.StepRetryConfig
}

// NewStepRetryResolver constructs a resolver from global and per-step configs.
func NewStepRetryResolver(global config.RetryConfig, perStep config.StepRetryConfig) *StepRetryResolver {
	return &StepRetryResolver{global: global, perStep: perStep}
}

// Resolve returns the RetryPolicy that should be used for the given step name.
// If a per-step override exists it takes precedence over the global config.
func (r *StepRetryResolver) Resolve(stepName string) RetryPolicy {
	base := config.RetryOverride{
		MaxAttempts: r.global.MaxAttempts,
		DelayMs:     int(r.global.Delay / time.Millisecond),
	}
	eff := r.perStep.LookupOverride(stepName, base)
	return RetryPolicy{
		MaxAttempts: eff.MaxAttempts,
		Delay:       time.Duration(eff.DelayMs) * time.Millisecond,
	}
}
