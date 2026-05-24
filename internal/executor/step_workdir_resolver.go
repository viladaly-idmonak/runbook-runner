package executor

import (
	"os"

	"github.com/your-org/runbook-runner/internal/config"
	"github.com/your-org/runbook-runner/internal/parser"
)

// StepWorkdirResolver resolves the working directory for a given step,
// applying per-step overrides from config before falling back to the
// process working directory.
type StepWorkdirResolver struct {
	cfg     config.StepWorkdirConfig
	procCwd string
}

// NewStepWorkdirResolver creates a resolver using the provided config.
// procCwd is used as the ultimate fallback; pass "" to use os.Getwd().
func NewStepWorkdirResolver(cfg config.StepWorkdirConfig, procCwd string) (*StepWorkdirResolver, error) {
	if procCwd == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		procCwd = cwd
	}
	return &StepWorkdirResolver{cfg: cfg, procCwd: procCwd}, nil
}

// Resolve returns the working directory for step. The priority is:
//  1. Per-step override from config
//  2. Config default (if non-empty)
//  3. Process working directory
func (r *StepWorkdirResolver) Resolve(step parser.Step) string {
	if dir := r.cfg.DirForStep(step.Name); dir != "" {
		return dir
	}
	return r.procCwd
}

// ApplyToEnv returns a copy of env with DIR set to the resolved directory
// for the step. This is a convenience helper for callers that build
// os/exec.Cmd environments manually.
func (r *StepWorkdirResolver) ApplyToEnv(step parser.Step, env []string) []string {
	dir := r.Resolve(step)
	result := make([]string, 0, len(env)+1)
	for _, e := range env {
		result = append(result, e)
	}
	result = append(result, "STEP_WORKDIR="+dir)
	return result
}
