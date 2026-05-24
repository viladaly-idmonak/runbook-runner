package executor

import (
	"context"
	"testing"
	"time"

	"github.com/user/runbook-runner/internal/config"
)

func TestRateLimiter_DisabledWaitIsNoop(t *testing.T) {
	cfg := config.RateConfig{Enabled: false}
	rl := NewRateLimiter(cfg)
	ctx := context.Background()
	start := time.Now()
	for i := 0; i < 20; i++ {
		if err := rl.Wait(ctx); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if elapsed := time.Since(start); elapsed > 50*time.Millisecond {
		t.Errorf("disabled limiter should not delay, elapsed %v", elapsed)
	}
}

func TestRateLimiter_EnabledAllowsBurst(t *testing.T) {
	cfg := config.RateConfig{
		Enabled:      true,
		MaxPerMinute: 600, // 10/s
		Burst:        5,
	}
	rl := NewRateLimiter(cfg)
	ctx := context.Background()
	// Burst of 5 should complete near-instantly
	start := time.Now()
	for i := 0; i < 5; i++ {
		if err := rl.Wait(ctx); err != nil {
			t.Fatalf("unexpected error on call %d: %v", i, err)
		}
	}
	if elapsed := time.Since(start); elapsed > 200*time.Millisecond {
		t.Errorf("burst should complete quickly, elapsed %v", elapsed)
	}
}

func TestRateLimiter_MinStepDelayIsRespected(t *testing.T) {
	delay := 80 * time.Millisecond
	cfg := config.RateConfig{
		Enabled:      true,
		MaxPerMinute: 3600, // 60/s — effectively no token limit
		Burst:        60,
		MinStepDelay: delay,
	}
	rl := NewRateLimiter(cfg)
	ctx := context.Background()
	start := time.Now()
	for i := 0; i < 2; i++ {
		if err := rl.Wait(ctx); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if elapsed := time.Since(start); elapsed < delay {
		t.Errorf("expected at least %v delay, got %v", delay, elapsed)
	}
}

func TestRateLimiter_ContextCancellation(t *testing.T) {
	cfg := config.RateConfig{
		Enabled:      true,
		MaxPerMinute: 1, // 1 per 60s — very slow
		Burst:        1,
	}
	rl := NewRateLimiter(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	// First call consumes the burst token.
	_ = rl.Wait(ctx)
	// Second call should block and then be cancelled.
	err := rl.Wait(ctx)
	if err == nil {
		t.Error("expected context cancellation error")
	}
}
