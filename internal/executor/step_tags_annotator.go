package executor

import (
	"runbook-runner/internal/config"
	"runbook-runner/internal/parser"
)

// StepTagsAnnotator attaches tag metadata to steps based on config.
type StepTagsAnnotator struct {
	cfg config.StepTagsConfig
}

// NewStepTagsAnnotator creates a StepTagsAnnotator from the given config.
func NewStepTagsAnnotator(cfg config.StepTagsConfig) *StepTagsAnnotator {
	return &StepTagsAnnotator{cfg: cfg}
}

// Annotate returns a copy of steps with Tags populated from config.
// If the config is disabled, the original slice is returned unchanged.
func (a *StepTagsAnnotator) Annotate(steps []parser.Step) []parser.Step {
	if !a.cfg.Enabled {
		return steps
	}
	annotated := make([]parser.Step, len(steps))
	for i, s := range steps {
		tags := config.TagsForStep(a.cfg, s.Name)
		if len(tags) > 0 {
			merged := make([]string, 0, len(s.Labels)+len(tags))
			merged = append(merged, s.Labels...)
			merged = append(merged, tags...)
			s.Labels = merged
		}
		annotated[i] = s
	}
	return annotated
}

// HasTag reports whether the given step has been annotated with a specific tag.
func (a *StepTagsAnnotator) HasTag(step parser.Step, tag string) bool {
	for _, l := range step.Labels {
		if l == tag {
			return true
		}
	}
	return false
}
