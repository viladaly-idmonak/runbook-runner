package executor

import (
	"github.com/your-org/runbook-runner/internal/config"
	"github.com/your-org/runbook-runner/internal/parser"
)

// LabelStepFilter wraps a slice of parsed steps and exposes only those
// that pass the label filter derived from the provided LabelConfig.
type LabelStepFilter struct {
	filter *config.LabelFilter
}

// NewLabelStepFilter constructs a LabelStepFilter from a LabelConfig.
func NewLabelStepFilter(cfg config.LabelConfig) *LabelStepFilter {
	return &LabelStepFilter{filter: config.NewLabelFilter(cfg)}
}

// Apply returns the subset of steps whose labels satisfy the filter.
// Steps that carry no labels are treated as having an empty label set.
func (f *LabelStepFilter) Apply(steps []parser.Step) []parser.Step {
	out := make([]parser.Step, 0, len(steps))
	for _, s := range steps {
		labels := config.StepLabels(s.Labels)
		if f.filter.Allow(labels) {
			out = append(out, s)
		}
	}
	return out
}
