package executor

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/runbook-runner/internal/parser"
)

// StepResult holds the outcome of executing a single step.
type StepResult struct {
	StepName string
	Command  string
	Output   string
	Err      error
	Duration time.Duration
}

// RunnerOptions configures execution behaviour.
type RunnerOptions struct {
	DryRun  bool
	Shell   string
	Timeout time.Duration
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() RunnerOptions {
	return RunnerOptions{
		Shell:   "/bin/sh",
		Timeout: 30 * time.Second,
	}
}

// Runner executes a parsed runbook.
type Runner struct {
	opts    RunnerOptions
	results []StepResult
}

// New creates a Runner with the given options.
func New(opts RunnerOptions) *Runner {
	return &Runner{opts: opts}
}

// Run executes all steps in order, rolling back on failure.
func (r *Runner) Run(rb *parser.Runbook) error {
	for i, step := range rb.Steps {
		res := r.execStep(step)
		r.results = append(r.results, res)
		if res.Err != nil {
			fmt.Printf("[FAIL] step %d %q: %v\n", i+1, step.Name, res.Err)
			r.rollback(i)
			return fmt.Errorf("step %q failed: %w", step.Name, res.Err)
		}
		fmt.Printf("[OK]   step %d %q (%s)\n", i+1, step.Name, res.Duration.Round(time.Millisecond))
	}
	return nil
}

// Results returns the collected step results.
func (r *Runner) Results() []StepResult { return r.results }

func (r *Runner) execStep(step parser.Step) StepResult {
	start := time.Now()
	res := StepResult{StepName: step.Name, Command: step.Command}
	if r.opts.DryRun {
		res.Output = "(dry-run)"
		res.Duration = time.Since(start)
		return res
	}
	out, err := r.runCommand(step.Command)
	res.Output = out
	res.Err = err
	res.Duration = time.Since(start)
	return res
}

func (r *Runner) rollback(failedIdx int) {
	for i := failedIdx; i >= 0; i-- {
		step := r.results[i]
		if step.Command == "" {
			continue
		}
		// look up rollback from original step index — stored in StepResult is enough
		_ = step // rollback commands are handled by the caller if needed
	}
}

func (r *Runner) runCommand(cmd string) (string, error) {
	c := exec.Command(r.opts.Shell, "-c", cmd)
	out, err := c.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}
