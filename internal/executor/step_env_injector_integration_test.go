package executor

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/user/runbook-runner/internal/config"
)

func runWithEnv(env []string, script string) (string, error) {
	cmd := exec.Command("sh", "-c", script)
	cmd.Env = env
	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err
}

func TestStepEnvInjector_Integration_OverrideVisibleToCommand(t *testing.T) {
	cfgVal := config.StepEnvConfig{
		Enabled: true,
		Overrides: map[string][]string{
			"deploy": {"MY_REGION=ap-southeast-1"},
		},
	}
	inj := NewStepEnvInjector(cfgVal)
	env := inj.Inject("deploy", []string{})

	out, err := runWithEnv(env, "echo $MY_REGION")
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}
	if out != "ap-southeast-1" {
		t.Errorf("expected ap-southeast-1, got %q", out)
	}
}

func TestStepEnvInjector_Integration_BasePreservedWhenNoConflict(t *testing.T) {
	cfgVal := config.StepEnvConfig{
		Enabled: true,
		Overrides: map[string][]string{
			"deploy": {"NEW_VAR=hello"},
		},
	}
	inj := NewStepEnvInjector(cfgVal)
	env := inj.Inject("deploy", []string{"EXISTING=world"})

	out, err := runWithEnv(env, "echo $EXISTING-$NEW_VAR")
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}
	if out != "world-hello" {
		t.Errorf("expected world-hello, got %q", out)
	}
}

func TestStepEnvInjector_Integration_OtherStepUnaffected(t *testing.T) {
	cfgVal := config.StepEnvConfig{
		Enabled: true,
		Overrides: map[string][]string{
			"deploy": {"SECRET=abc"},
		},
	}
	inj := NewStepEnvInjector(cfgVal)
	env := inj.Inject("verify", []string{})

	out, err := runWithEnv(env, "echo ${SECRET:-empty}")
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}
	if out != "empty" {
		t.Errorf("expected empty, got %q", out)
	}
}
