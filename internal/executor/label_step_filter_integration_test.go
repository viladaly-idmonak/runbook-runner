package executor_test

import (
	"testing"

	"github.com/your-org/runbook-runner/internal/config"
	"github.com/your-org/runbook-runner/internal/executor"
	"github.com/your-org/runbook-runner/internal/parser"
)

func TestLabelStepFilter_Integration_OnlyProdStepsRun(t *testing.T) {
	cfg := config.DefaultLabelConfig()
	cfg.Include = map[string]string{"env": "prod"}

	filter := executor.NewLabelStepFilter(cfg)

	allSteps := []parser.Step{
		{Name: "deploy", Command: "echo deploy", Labels: map[string]string{"env": "prod"}},
		{Name: "smoke-test", Command: "echo smoke", Labels: map[string]string{"env": "staging"}},
		{Name: "notify", Command: "echo notify", Labels: map[string]string{"env": "prod", "team": "ops"}},
	}

	filtered := filter.Apply(allSteps)

	if len(filtered) != 2 {
		t.Fatalf("expected 2 prod steps, got %d", len(filtered))
	}
	for _, s := range filtered {
		if s.Labels["env"] != "prod" {
			t.Errorf("non-prod step leaked through filter: %s", s.Name)
		}
	}
}

func TestLabelStepFilter_Integration_ExcludeSkipsStep(t *testing.T) {
	cfg := config.DefaultLabelConfig()
	cfg.Exclude = map[string]string{"dangerous": "true"}

	filter := executor.NewLabelStepFilter(cfg)

	allSteps := []parser.Step{
		{Name: "safe", Command: "echo safe", Labels: map[string]string{}},
		{Name: "risky", Command: "rm -rf /tmp/test", Labels: map[string]string{"dangerous": "true"}},
	}

	filtered := filter.Apply(allSteps)

	if len(filtered) != 1 {
		t.Fatalf("expected 1 step after exclusion, got %d", len(filtered))
	}
	if filtered[0].Name != "safe" {
		t.Errorf("expected 'safe' step, got %q", filtered[0].Name)
	}
}
