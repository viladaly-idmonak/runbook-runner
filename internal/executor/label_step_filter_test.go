package executor

import (
	"testing"

	"github.com/your-org/runbook-runner/internal/config"
	"github.com/your-org/runbook-runner/internal/parser"
)

func makeSteps(labelSets ...map[string]string) []parser.Step {
	steps := make([]parser.Step, len(labelSets))
	for i, ls := range labelSets {
		steps[i] = parser.Step{
			Name:   "step",
			Labels: ls,
		}
	}
	return steps
}

func TestLabelStepFilter_NoRules_ReturnsAll(t *testing.T) {
	f := NewLabelStepFilter(config.DefaultLabelConfig())
	steps := makeSteps(map[string]string{"env": "prod"}, map[string]string{})
	got := f.Apply(steps)
	if len(got) != 2 {
		t.Fatalf("expected 2 steps, got %d", len(got))
	}
}

func TestLabelStepFilter_IncludeFiltersCorrectly(t *testing.T) {
	cfg := config.DefaultLabelConfig()
	cfg.Include = map[string]string{"env": "prod"}
	f := NewLabelStepFilter(cfg)
	steps := makeSteps(
		map[string]string{"env": "prod"},
		map[string]string{"env": "staging"},
	)
	got := f.Apply(steps)
	if len(got) != 1 {
		t.Fatalf("expected 1 step, got %d", len(got))
	}
	if got[0].Labels["env"] != "prod" {
		t.Errorf("unexpected step label: %v", got[0].Labels)
	}
}

func TestLabelStepFilter_ExcludeFiltersCorrectly(t *testing.T) {
	cfg := config.DefaultLabelConfig()
	cfg.Exclude = map[string]string{"skip": "true"}
	f := NewLabelStepFilter(cfg)
	steps := makeSteps(
		map[string]string{"skip": "true"},
		map[string]string{"skip": "false"},
		map[string]string{},
	)
	got := f.Apply(steps)
	if len(got) != 2 {
		t.Fatalf("expected 2 steps, got %d", len(got))
	}
}

func TestLabelStepFilter_EmptySteps_ReturnsEmpty(t *testing.T) {
	f := NewLabelStepFilter(config.DefaultLabelConfig())
	got := f.Apply([]parser.Step{})
	if len(got) != 0 {
		t.Errorf("expected empty result, got %d steps", len(got))
	}
}
