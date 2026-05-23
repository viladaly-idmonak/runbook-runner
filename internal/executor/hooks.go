package executor

import (
	"fmt"
	"os/exec"
	"strings"
)

// HookType identifies when a hook runs relative to step execution.
type HookType string

const (
	HookPreStep  HookType = "pre-step"
	HookPostStep HookType = "post-step"
	HookOnError  HookType = "on-error"
)

// Hook represents a shell command to run at a specific lifecycle point.
type Hook struct {
	Type    HookType
	Command string
}

// HookRunner executes lifecycle hooks with access to step context.
type HookRunner struct {
	shell   string
	hooks   []Hook
	verbose bool
}

// NewHookRunner creates a HookRunner with the given shell and registered hooks.
func NewHookRunner(shell string, hooks []Hook, verbose bool) *HookRunner {
	return &HookRunner{shell: shell, hooks: hooks, verbose: verbose}
}

// Run executes all hooks of the given type, injecting step metadata as env vars.
func (h *HookRunner) Run(hookType HookType, stepIndex int, stepName string) error {
	for _, hook := range h.hooks {
		if hook.Type != hookType {
			continue
		}
		if err := h.execute(hook.Command, stepIndex, stepName); err != nil {
			return fmt.Errorf("hook %s failed for step %q: %w", hookType, stepName, err)
		}
	}
	return nil
}

func (h *HookRunner) execute(command string, stepIndex int, stepName string) error {
	cmd := exec.Command(h.shell, "-c", command)
	cmd.Env = append(cmd.Environ(),
		fmt.Sprintf("RR_STEP_INDEX=%d", stepIndex),
		fmt.Sprintf("RR_STEP_NAME=%s", stepName),
	)
	out, err := cmd.CombinedOutput()
	if h.verbose && len(strings.TrimSpace(string(out))) > 0 {
		fmt.Printf("[hook] %s\n", strings.TrimSpace(string(out)))
	}
	return err
}
