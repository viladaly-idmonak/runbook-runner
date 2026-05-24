package executor

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/example/runbook-runner/internal/config"
)

// ConditionEvaluator decides whether a step should be executed.
type ConditionEvaluator struct {
	cfg   config.ConditionConfig
	shell string
}

// NewConditionEvaluator creates a ConditionEvaluator using the provided config.
func NewConditionEvaluator(cfg config.ConditionConfig, shell string) *ConditionEvaluator {
	if shell == "" {
		shell = "/bin/sh"
	}
	return &ConditionEvaluator{cfg: cfg, shell: shell}
}

// Evaluate returns true when the step should run.
// expr is the raw condition expression; mode overrides the default when non-empty.
// If conditions are disabled the method always returns true.
func (e *ConditionEvaluator) Evaluate(expr string, mode config.ConditionMode) (bool, error) {
	if !e.cfg.Enabled || expr == "" {
		return true, nil
	}
	if mode == "" {
		mode = e.cfg.DefaultMode
	}
	switch mode {
	case config.ConditionModeEnvSet:
		return e.evalEnvSet(expr)
	case config.ConditionModeShell:
		return e.evalShell(expr)
	default:
		return false, fmt.Errorf("condition: unsupported mode %q", mode)
	}
}

// evalEnvSet returns true when the named environment variable is non-empty.
func (e *ConditionEvaluator) evalEnvSet(varName string) (bool, error) {
	varName = strings.TrimSpace(varName)
	val, ok := os.LookupEnv(varName)
	if !ok || val == "" {
		if e.cfg.SkipOnConditionFailure {
			return false, nil
		}
		return false, fmt.Errorf("condition: env var %q is not set", varName)
	}
	return true, nil
}

// evalShell runs expr as a shell command and returns true on exit code 0.
func (e *ConditionEvaluator) evalShell(expr string) (bool, error) {
	cmd := exec.Command(e.shell, "-c", expr) //nolint:gosec
	if err := cmd.Run(); err != nil {
		if e.cfg.SkipOnConditionFailure {
			return false, nil
		}
		return false, fmt.Errorf("condition: shell expression failed: %w", err)
	}
	return true, nil
}
