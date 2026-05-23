package executor

import (
	"fmt"

	"github.com/runbook-runner/internal/parser"
)

// RollbackResult captures the outcome of a single rollback command.
type RollbackResult struct {
	StepName string
	Command  string
	Output   string
	Err      error
}

// Rollback executes rollback commands for steps[0..failedIdx] in reverse order.
// Steps without a rollback command are silently skipped.
func (r *Runner) Rollback(steps []parser.Step, failedIdx int) []RollbackResult {
	var results []RollbackResult
	for i := failedIdx; i >= 0; i-- {
		step := steps[i]
		if step.Rollback == "" {
			continue
		}
		res := RollbackResult{StepName: step.Name, Command: step.Rollback}
		if r.opts.DryRun {
			res.Output = "(dry-run rollback)"
			fmt.Printf("[ROLLBACK-DRY] step %d %q: %s\n", i+1, step.Name, step.Rollback)
		} else {
			out, err := r.runCommand(step.Rollback)
			res.Output = out
			res.Err = err
			if err != nil {
				fmt.Printf("[ROLLBACK-FAIL] step %d %q: %v\n", i+1, step.Name, err)
			} else {
				fmt.Printf("[ROLLBACK-OK]  step %d %q\n", i+1, step.Name)
			}
		}
		results = append(results, res)
	}
	return results
}
