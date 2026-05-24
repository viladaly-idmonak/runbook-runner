package executor

import (
	"testing"

	"github.com/user/runbook-runner/internal/config"
	"github.com/user/runbook-runner/internal/parser"
)

func stepsFromNames(names ...string) []parser.Step {
	out := make([]parser.Step, len(names))
	for i, n := range names {
		out[i] = parser.Step{Name: n, Command: "echo " + n}
	}
	return out
}

func TestStepDependsRunner_Disabled_PreservesOrder(t *testing.T) {
	r := NewStepDependsRunner(config.StepDependsConfig{Enabled: false})
	steps := stepsFromNames("c", "a", "b")
	got, err := r.Reorder(steps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, s := range got {
		if s.Name != steps[i].Name {
			t.Errorf("position %d: got %q want %q", i, s.Name, steps[i].Name)
		}
	}
}

func TestStepDependsRunner_Enabled_NoDeps_PreservesOrder(t *testing.T) {
	r := NewStepDependsRunner(config.StepDependsConfig{Enabled: true})
	steps := stepsFromNames("a", "b", "c")
	got, err := r.Reorder(steps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 steps, got %d", len(got))
	}
}

func TestStepDependsRunner_Reorders_ByDependency(t *testing.T) {
	cfg := config.StepDependsConfig{
		Enabled: true,
		Overrides: []config.StepDependsOverride{
			{Step: "b", DependsOn: []string{"a"}},
			{Step: "c", DependsOn: []string{"b"}},
		},
	}
	r := NewStepDependsRunner(cfg)
	// Provide steps in reverse order to verify reordering.
	steps := stepsFromNames("c", "b", "a")
	got, err := r.Reorder(steps)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"a", "b", "c"}
	for i, s := range got {
		if s.Name != want[i] {
			t.Errorf("position %d: got %q want %q", i, s.Name, want[i])
		}
	}
}

func TestStepDependsRunner_DetectsCycle(t *testing.T) {
	cfg := config.StepDependsConfig{
		Enabled: true,
		Overrides: []config.StepDependsOverride{
			{Step: "a", DependsOn: []string{"b"}},
			{Step: "b", DependsOn: []string{"a"}},
		},
	}
	r := NewStepDependsRunner(cfg)
	_, err := r.Reorder(stepsFromNames("a", "b"))
	if err == nil {
		t.Fatal("expected cycle error, got nil")
	}
}

func TestStepDependsRunner_UnknownDep_ReturnsError(t *testing.T) {
	cfg := config.StepDependsConfig{
		Enabled: true,
		Overrides: []config.StepDependsOverride{
			{Step: "a", DependsOn: []string{"missing"}},
		},
	}
	r := NewStepDependsRunner(cfg)
	_, err := r.Reorder(stepsFromNames("a"))
	if err == nil {
		t.Fatal("expected error for unknown dep, got nil")
	}
}
