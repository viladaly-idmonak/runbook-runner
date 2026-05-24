package config

import "fmt"

// SkipConfig controls which steps should be unconditionally skipped during execution.
type SkipConfig struct {
	// StepNames is a list of step titles to skip by exact match.
	StepNames []string `yaml:"step_names"`

	// StepIndices is a list of 0-based step indices to skip.
	StepIndices []int `yaml:"step_indices"`

	// OnlyWhen optionally skips steps when a shell expression evaluates to a
	// non-zero exit code (i.e. skip if the condition is NOT met).
	OnlyWhen string `yaml:"only_when"`
}

// DefaultSkipConfig returns a SkipConfig with no skips configured.
func DefaultSkipConfig() SkipConfig {
	return SkipConfig{
		StepNames:   []string{},
		StepIndices: []int{},
		OnlyWhen:    "",
	}
}

// ValidateSkip returns an error if the SkipConfig contains invalid values.
func ValidateSkip(c SkipConfig) error {
	seen := make(map[string]struct{}, len(c.StepNames))
	for _, name := range c.StepNames {
		if name == "" {
			return fmt.Errorf("skip: step_names must not contain empty strings")
		}
		if _, dup := seen[name]; dup {
			return fmt.Errorf("skip: duplicate step name %q", name)
		}
		seen[name] = struct{}{}
	}

	seenIdx := make(map[int]struct{}, len(c.StepIndices))
	for _, idx := range c.StepIndices {
		if idx < 0 {
			return fmt.Errorf("skip: step_indices must not contain negative values, got %d", idx)
		}
		if _, dup := seenIdx[idx]; dup {
			return fmt.Errorf("skip: duplicate step index %d", idx)
		}
		seenIdx[idx] = struct{}{}
	}

	return nil
}

// ShouldSkip reports whether the step at the given index with the given title
// should be skipped according to c.
func (c SkipConfig) ShouldSkip(index int, title string) bool {
	for _, name := range c.StepNames {
		if name == title {
			return true
		}
	}
	for _, idx := range c.StepIndices {
		if idx == index {
			return true
		}
	}
	return false
}
