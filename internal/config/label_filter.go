package config

// StepLabels represents the labels attached to a runbook step.
type StepLabels map[string]string

// LabelFilter evaluates whether a step should be executed based on
// the include/exclude rules defined in LabelConfig.
type LabelFilter struct {
	cfg LabelConfig
}

// NewLabelFilter creates a LabelFilter backed by the given LabelConfig.
func NewLabelFilter(cfg LabelConfig) *LabelFilter {
	return &LabelFilter{cfg: cfg}
}

// Allow returns true when the step with the given labels should run.
// A step is allowed when:
//  1. It matches every include label (or no include labels are defined).
//  2. It does NOT match any exclude label.
func (f *LabelFilter) Allow(labels StepLabels) bool {
	for k, v := range f.cfg.Include {
		if labels[k] != v {
			return false
		}
	}
	for k, v := range f.cfg.Exclude {
		if labels[k] == v {
			return false
		}
	}
	return true
}
