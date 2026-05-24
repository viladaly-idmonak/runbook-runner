package executor

import (
	"fmt"

	"github.com/user/runbook-runner/internal/config"
	"github.com/user/runbook-runner/internal/parser"
)

// StepDependsRunner reorders and optionally skips steps based on declared
// dependencies using a topological sort.
type StepDependsRunner struct {
	cfg config.StepDependsConfig
}

// NewStepDependsRunner returns a StepDependsRunner for the given config.
func NewStepDependsRunner(cfg config.StepDependsConfig) *StepDependsRunner {
	return &StepDependsRunner{cfg: cfg}
}

// Reorder returns steps in dependency-safe execution order. When the feature
// is disabled the original slice is returned unchanged.
func (r *StepDependsRunner) Reorder(steps []parser.Step) ([]parser.Step, error) {
	if !r.cfg.Enabled {
		return steps, nil
	}

	names := make([]string, len(steps))
	index := make(map[string]parser.Step, len(steps))
	for i, s := range steps {
		names[i] = s.Name
		if _, dup := index[s.Name]; dup {
			return nil, fmt.Errorf("duplicate step name %q", s.Name)
		}
		index[s.Name] = s
	}

	// Build dependency map from config overrides.
	deps := make(map[string][]string)
	for _, o := range r.cfg.Overrides {
		deps[o.Step] = o.DependsOn
	}

	ordered, err := config.DependencyOrder(names, deps)
	if err != nil {
		return nil, err
	}

	result := make([]parser.Step, 0, len(ordered))
	for _, name := range ordered {
		result = append(result, index[name])
	}
	return result, nil
}
