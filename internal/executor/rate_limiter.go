package executor

import (
	"context"
	"time"

	"golang.org/x/time/rate"

	"github.com/user/runbook-runner/internal/config"
)

// RateLimitedRunner wraps a RunFunc and enforces a token-bucket rate limit
// plus an optional minimum delay between steps.
type RateLimitedRunner struct {
	limiter      *rate.Limiter
	minStepDelay time.Duration
	enabled      bool
}

// NewRateLimiter creates a RateLimitedRunner from the given RateConfig.
// If the config is disabled, Wait becomes a no-op.
func NewRateLimiter(cfg config.RateConfig) *RateLimitedRunner {
	if !cfg.Enabled {
		return &RateLimitedRunner{enabled: false}
	}

	tokensPerSecond := rate.Limit(float64(cfg.MaxPerMinute) / 60.0)
	burst := cfg.Burst
	if burst == 0 {
		burst = 1
	}

	return &RateLimitedRunner{
		limiter:      rate.NewLimiter(tokensPerSecond, burst),
		minStepDelay: cfg.MinStepDelay,
		enabled:      true,
	}
}

// Wait blocks until the rate limiter permits the next step to proceed.
// It respects context cancellation and enforces MinStepDelay after acquiring
// the token.
func (r *RateLimitedRunner) Wait(ctx context.Context) error {
	if !r.enabled {
		return nil
	}
	if err := r.limiter.Wait(ctx); err != nil {
		return err
	}
	if r.minStepDelay > 0 {
		select {
		case <-time.After(r.minStepDelay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}
