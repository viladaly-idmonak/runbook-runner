package executor

import (
	"testing"

	"github.com/user/runbook-runner/internal/config"
)

func cfg(enabled bool, overrides map[string][]string) config.StepEnvConfig {
	return config.StepEnvConfig{Enabled: enabled, Overrides: overrides}
}

func TestStepEnvInjector_DisabledReturnsBase(t *testing.T) {
	inj := NewStepEnvInjector(cfg(false, map[string][]string{"s": {"K=V"}}))
	base := []string{"A=1"}
	got := inj.Inject("s", base)
	if len(got) != 1 || got[0] != "A=1" {
		t.Errorf("expected original base, got %v", got)
	}
}

func TestStepEnvInjector_NoOverrideForStep(t *testing.T) {
	inj := NewStepEnvInjector(cfg(true, map[string][]string{"other": {"K=V"}}))
	base := []string{"A=1"}
	got := inj.Inject("deploy", base)
	if len(got) != 1 || got[0] != "A=1" {
		t.Errorf("expected unchanged base, got %v", got)
	}
}

func TestStepEnvInjector_AddsNewKeys(t *testing.T) {
	inj := NewStepEnvInjector(cfg(true, map[string][]string{"deploy": {"REGION=us-east-1"}}))
	base := []string{"HOME=/root"}
	got := inj.Inject("deploy", base)
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d: %v", len(got), got)
	}
}

func TestStepEnvInjector_OverridesExistingKey(t *testing.T) {
	inj := NewStepEnvInjector(cfg(true, map[string][]string{"deploy": {"REGION=eu-west-1"}}))
	base := []string{"REGION=us-east-1", "HOME=/root"}
	got := inj.Inject("deploy", base)
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d: %v", len(got), got)
	}
	for _, e := range got {
		if e == "REGION=us-east-1" {
			t.Error("old REGION value should have been replaced")
		}
	}
}

func TestStepEnvInjector_PairsDisabled(t *testing.T) {
	inj := NewStepEnvInjector(cfg(false, map[string][]string{"s": {"K=V"}}))
	if p := inj.Pairs("s"); p != nil {
		t.Errorf("expected nil when disabled, got %v", p)
	}
}

func TestStepEnvInjector_PairsEnabled(t *testing.T) {
	inj := NewStepEnvInjector(cfg(true, map[string][]string{"s": {"K=V"}}))
	if p := inj.Pairs("s"); len(p) != 1 || p[0] != "K=V" {
		t.Errorf("unexpected pairs: %v", p)
	}
}
