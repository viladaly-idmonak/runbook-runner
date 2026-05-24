package config

import "fmt"

// OnFailureAction defines what the runner should do when a step fails.
type OnFailureAction string

const (
	OnFailureContinue  OnFailureAction = "continue"
	OnFailureStop      OnFailureAction = "stop"
	OnFailureRollback  OnFailureAction = "rollback"
)

var validOnFailureActions = map[OnFailureAction]struct{}{
	OnFailureContinue: {},
	OnFailureStop:     {},
	OnFailureRollback: {},
}

// StepOnFailureConfig controls per-step failure behaviour overrides.
type StepOnFailureConfig struct {
	Enabled   bool
	// Overrides maps step name to the desired OnFailureAction.
	Overrides map[string]OnFailureAction
	// Default is the fallback action when no per-step override exists.
	Default OnFailureAction
}

// DefaultStepOnFailureConfig returns a config that stops on first failure.
func DefaultStepOnFailureConfig() StepOnFailureConfig {
	return StepOnFailureConfig{
		Enabled:   false,
		Overrides: map[string]OnFailureAction{},
		Default:   OnFailureStop,
	}
}

// ValidateStepOnFailure returns an error if the config is invalid.
func ValidateStepOnFailure(c StepOnFailureConfig) error {
	if !c.Enabled {
		return nil
	}
	if _, ok := validOnFailureActions[c.Default]; !ok {
		return fmt.Errorf("step_on_failure: unknown default action %q", c.Default)
	}
	for step, action := range c.Overrides {
		if step == "" {
			return fmt.Errorf("step_on_failure: override has empty step name")
		}
		if _, ok := validOnFailureActions[action]; !ok {
			return fmt.Errorf("step_on_failure: step %q has unknown action %q", step, action)
		}
	}
	return nil
}

// ActionForStep returns the resolved OnFailureAction for the given step name.
func ActionForStep(c StepOnFailureConfig, stepName string) OnFailureAction {
	if !c.Enabled {
		return OnFailureStop
	}
	if action, ok := c.Overrides[stepName]; ok {
		return action
	}
	return c.Default
}
