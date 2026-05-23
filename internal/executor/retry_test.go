package executor

import (
	"errors"
	"testing"
	"time"
)

// stubRunner is a CommandRunner that returns pre-configured results in order.
type stubRunner struct {
	results []stubResult
	callIdx int
}

type stubResult struct {
	out string
	err error
}

func (s *stubRunner) Run(_ string) (string, error) {
	if s.callIdx >= len(s.results) {
		return "", errors.New("no more stub results")
	}
	r := s.results[s.callIdx]
	s.callIdx++
	return r.out, r.err
}

func TestRetryRunner_SuccessOnFirstAttempt(t *testing.T) {
	stub := &stubRunner{results: []stubResult{{out: "ok", err: nil}}}
	rr := NewRetryRunner(stub, RetryPolicy{MaxAttempts: 3, Delay: 0})
	out, err := rr.Run("echo ok")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "ok" {
		t.Errorf("expected 'ok', got %q", out)
	}
	if stub.callIdx != 1 {
		t.Errorf("expected 1 call, got %d", stub.callIdx)
	}
}

func TestRetryRunner_SuccessOnSecondAttempt(t *testing.T) {
	stub := &stubRunner{results: []stubResult{
		{out: "", err: errors.New("fail")},
		{out: "recovered", err: nil},
	}}
	rr := NewRetryRunner(stub, RetryPolicy{MaxAttempts: 3, Delay: 0})
	out, err := rr.Run("cmd")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "recovered" {
		t.Errorf("expected 'recovered', got %q", out)
	}
}

func TestRetryRunner_AllAttemptsExhausted(t *testing.T) {
	stub := &stubRunner{results: []stubResult{
		{err: errors.New("fail")},
		{err: errors.New("fail")},
		{err: errors.New("fail")},
	}}
	rr := NewRetryRunner(stub, RetryPolicy{MaxAttempts: 3, Delay: 0})
	_, err := rr.Run("cmd")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if stub.callIdx != 3 {
		t.Errorf("expected 3 calls, got %d", stub.callIdx)
	}
}

func TestRetryRunner_ZeroMaxAttemptsDefaultsToOne(t *testing.T) {
	stub := &stubRunner{results: []stubResult{{out: "hi", err: nil}}}
	rr := NewRetryRunner(stub, RetryPolicy{MaxAttempts: 0})
	_, err := rr.Run("cmd")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stub.callIdx != 1 {
		t.Errorf("expected 1 call, got %d", stub.callIdx)
	}
}

func TestRetryRunner_DelayBetweenAttempts(t *testing.T) {
	stub := &stubRunner{results: []stubResult{
		{err: errors.New("fail")},
		{out: "ok", err: nil},
	}}
	rr := NewRetryRunner(stub, RetryPolicy{MaxAttempts: 2, Delay: 10 * time.Millisecond})
	start := time.Now()
	_, err := rr.Run("cmd")
	elapsed := time.Since(start)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if elapsed < 10*time.Millisecond {
		t.Errorf("expected delay >= 10ms, got %v", elapsed)
	}
}
