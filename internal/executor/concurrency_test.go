package executor

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func echoRun(_ context.Context, cmd string) (string, error) {
	return "out:" + cmd, nil
}

func TestConcurrentRunner_ResultsInOrder(t *testing.T) {
	r := NewConcurrentRunner(3, echoRun)
	cmds := []string{"a", "b", "c", "d"}
	results := r.RunAll(context.Background(), cmds)

	if len(results) != len(cmds) {
		t.Fatalf("expected %d results, got %d", len(cmds), len(results))
	}
	for i, res := range results {
		if res.Index != i {
			t.Errorf("result[%d].Index = %d", i, res.Index)
		}
		want := "out:" + cmds[i]
		if res.Output != want {
			t.Errorf("result[%d].Output = %q, want %q", i, res.Output, want)
		}
		if res.Err != nil {
			t.Errorf("result[%d].Err = %v", i, res.Err)
		}
	}
}

func TestConcurrentRunner_WorkerLimit(t *testing.T) {
	var active int64
	var peak int64

	slowRun := func(_ context.Context, cmd string) (string, error) {
		cur := atomic.AddInt64(&active, 1)
		for {
			p := atomic.LoadInt64(&peak)
			if cur <= p || atomic.CompareAndSwapInt64(&peak, p, cur) {
				break
			}
		}
		time.Sleep(10 * time.Millisecond)
		atomic.AddInt64(&active, -1)
		return cmd, nil
	}

	const workers = 2
	r := NewConcurrentRunner(workers, slowRun)
	cmds := make([]string, 8)
	for i := range cmds {
		cmds[i] = fmt.Sprintf("cmd%d", i)
	}
	r.RunAll(context.Background(), cmds)

	if peak > workers {
		t.Errorf("peak concurrency %d exceeded worker limit %d", peak, workers)
	}
}

func TestConcurrentRunner_PropagatesErrors(t *testing.T) {
	errBoom := errors.New("boom")
	failRun := func(_ context.Context, cmd string) (string, error) {
		if cmd == "bad" {
			return "", errBoom
		}
		return cmd, nil
	}

	r := NewConcurrentRunner(2, failRun)
	results := r.RunAll(context.Background(), []string{"ok", "bad", "ok2"})

	if results[1].Err == nil {
		t.Fatal("expected error for 'bad' command")
	}
	if results[0].Err != nil || results[2].Err != nil {
		t.Error("unexpected errors on good commands")
	}
}

func TestConcurrentRunner_DefaultsToOneWorker(t *testing.T) {
	r := NewConcurrentRunner(0, echoRun)
	if r.workers != 1 {
		t.Errorf("expected workers=1, got %d", r.workers)
	}
}

func TestFirstError_NilWhenAllOK(t *testing.T) {
	results := []StepResult{{Index: 0, Err: nil}, {Index: 1, Err: nil}}
	if err := FirstError(results); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestFirstError_ReturnsFirst(t *testing.T) {
	results := []StepResult{
		{Index: 0, Cmd: "ok", Err: nil},
		{Index: 1, Cmd: "fail", Err: errors.New("oops")},
	}
	err := FirstError(results)
	if err == nil {
		t.Fatal("expected non-nil error")
	}
}
