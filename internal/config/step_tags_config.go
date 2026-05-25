package config

import "fmt"

// StepTagsConfig controls per-step tag annotations used for filtering and reporting.
type StepTagsConfig struct {
	Enabled bool                `yaml:"enabled"`
	StepTags map[string][]string `yaml:"step_tags"`
}

// DefaultStepTagsConfig returns a StepTagsConfig with safe defaults.
func DefaultStepTagsConfig() StepTagsConfig {
	return StepTagsConfig{
		Enabled:  false,
		StepTags: map[string][]string{},
	}
}

// ValidateStepTags returns an error if the config is invalid.
func ValidateStepTags(c StepTagsConfig) error {
	if !c.Enabled {
		return nil
	}
	for stepName, tags := range c.StepTags {
		if stepName == "" {
			return fmt.Errorf("step_tags: step name must not be empty")
		}
		for i, tag := range tags {
			if tag == "" {
				return fmt.Errorf("step_tags: tag at index %d for step %q must not be empty", i, stepName)
			}
		}
	}
	return nil
}

// TagsForStep returns the tags associated with the given step name.
// Returns nil if the step has no tags or the config is disabled.
func TagsForStep(c StepTagsConfig, stepName string) []string {
	if !c.Enabled {
		return nil
	}
	return c.StepTags[stepName]
}
