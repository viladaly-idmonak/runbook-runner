package executor

import (
	"context"
	"fmt"
	"sync"
)

// ConcurrentRunner executes multiple independent steps in parallel up to a
// configurable worker limit and collects their results.
type ConcurrentRunner struct {
	workers int
	run     func(ctx context.Context, cmd string) (string, error)
}

// StepResult holds the outcome of a single concurrent step execution.
type StepResult struct {
	Index  int
	Cmd    string
	Output string
	Err    error
}

// NewConcurrentRunner returns a ConcurrentRunner that uses at most workers
// goroutines. If workers < 1 it defaults to 1.
func NewConcurrentRunner(workers int, run func(ctx context.Context, cmd string) (string, error)) *ConcurrentRunner {
	if workers < 1 {
		workers = 1
	}
	return &ConcurrentRunner{workers: workers, run: run}
}

// RunAll executes each command in cmds concurrently, respecting the worker
// limit, and returns one StepResult per command in the original order.
func (c *ConcurrentRunner) RunAll(ctx context.Context, cmds []string) []StepResult {
	results := make([]StepResult, len(cmds))
	work := make(chan int, len(cmds))

	for i := range cmds {
		work <- i
	}
	close(work)

	var wg sync.WaitGroup
	for w := 0; w < c.workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range work {
				out, err := c.run(ctx, cmds[idx])
				results[idx] = StepResult{
					Index:  idx,
					Cmd:    cmds[idx],
					Output: out,
					Err:    err,
				}
			}
		}()
	}
	wg.Wait()
	return results
}

// FirstError returns the first non-nil error from results, or nil.
func FirstError(results []StepResult) error {
	for _, r := range results {
		if r.Err != nil {
			return fmt.Errorf("step %d (%q): %w", r.Index, r.Cmd, r.Err)
		}
	}
	return nil
}
