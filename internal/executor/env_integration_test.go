package executor_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/user/runbook-runner/internal/config"
	"github.com/user/runbook-runner/internal/executor"
)

// TestEnvResolver_Integration_CommandSeesExtra verifies that a real subprocess
// launched with the resolved environment can read a variable injected via Extra.
func TestEnvResolver_Integration_CommandSeesExtra(t *testing.T) {
	cfg := config.DefaultEnvConfig()
	cfg.InheritAll = false
	cfg.Extra = map[string]string{"GREETING": "hello-world"}

	r := executor.NewEnvResolver(cfg)
	env := r.Resolve()

	cmd := exec.Command("sh", "-c", "echo $GREETING")
	cmd.Env = env
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}
	got := strings.TrimSpace(string(out))
	if got != "hello-world" {
		t.Errorf("expected 'hello-world', got %q", got)
	}
}

// TestEnvResolver_Integration_PassthroughForwardedToChild confirms that a
// variable listed in Passthrough reaches the child process.
func TestEnvResolver_Integration_PassthroughForwardedToChild(t *testing.T) {
	os.Setenv("_RR_INTEG", "forwarded")
	t.Cleanup(func() { os.Unsetenv("_RR_INTEG") })

	cfg := config.DefaultEnvConfig()
	cfg.InheritAll = false
	cfg.Passthrough = []string{"_RR_INTEG"}

	r := executor.NewEnvResolver(cfg)
	env := r.Resolve()

	cmd := exec.Command("sh", "-c", "echo $_RR_INTEG")
	cmd.Env = env
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}
	got := strings.TrimSpace(string(out))
	if got != "forwarded" {
		t.Errorf("expected 'forwarded', got %q", got)
	}
}
