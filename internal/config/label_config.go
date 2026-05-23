package config

import "fmt"

// LabelConfig controls step filtering by label.
type LabelConfig struct {
	// IncludeLabels, when non-empty, restricts execution to steps that carry
	// at least one of the listed labels.
	IncludeLabels []string `yaml:"include_labels"`

	// ExcludeLabels skips any step whose label list intersects with this set.
	ExcludeLabels []string `yaml:"exclude_labels"`
}

// DefaultLabelConfig returns a LabelConfig with no filtering applied.
func DefaultLabelConfig() LabelConfig {
	return LabelConfig{
		IncludeLabels: []string{},
		ExcludeLabels: []string{},
	}
}

// ValidateLabels returns an error if the LabelConfig is logically inconsistent.
func ValidateLabels(c LabelConfig) error {
	seen := make(map[string]struct{}, len(c.IncludeLabels))
	for _, l := range c.IncludeLabels {
		if l == "" {
			return fmt.Errorf("label config: include_labels contains an empty string")
		}
		seen[l] = struct{}{}
	}
	for _, l := range c.ExcludeLabels {
		if l == "" {
			return fmt.Errorf("label config: exclude_labels contains an empty string")
		}
		if _, conflict := seen[l]; conflict {
			return fmt.Errorf("label config: label %q appears in both include_labels and exclude_labels", l)
		}
	}
	return nil
}

// MatchesFilter reports whether a step with the given labels should be executed
// according to the provided LabelConfig.
func MatchesFilter(cfg LabelConfig, stepLabels []string) bool {
	labelSet := make(map[string]struct{}, len(stepLabels))
	for _, l := range stepLabels {
		labelSet[l] = struct{}{}
	}

	// Exclusion takes priority.
	for _, ex := range cfg.ExcludeLabels {
		if _, ok := labelSet[ex]; ok {
			return false
		}
	}

	// If no inclusion filter is set, accept everything.
	if len(cfg.IncludeLabels) == 0 {
		return true
	}

	for _, inc := range cfg.IncludeLabels {
		if _, ok := labelSet[inc]; ok {
			return true
		}
	}
	return false
}
