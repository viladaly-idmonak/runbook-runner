package reporter

import "time"

// StepStatus represents the outcome of a single runbook step.
type StepStatus string

const (
	StatusOK      StepStatus = "ok"
	StatusFail    StepStatus = "fail"
	StatusSkipped StepStatus = "skipped"
)

// String returns a human-readable label for the status.
func (s StepStatus) String() string {
	switch s {
	case StatusOK:
		return "OK"
	case StatusFail:
		return "FAIL"
	case StatusSkipped:
		return "SKIPPED"
	default:
		return string(s)
	}
}

// StepEntry records the result of executing one runbook step.
type StepEntry struct {
	// Name is the human-readable step title.
	Name string `json:"name"`

	// Status is the outcome of the step.
	Status StepStatus `json:"status"`

	// Output is the combined stdout/stderr of the step command.
	Output string `json:"output,omitempty"`

	// RollbackOutput is the combined output of the rollback command, if run.
	RollbackOutput string `json:"rollback_output,omitempty"`

	// Duration is how long the step took to execute.
	Duration time.Duration `json:"duration_ms"`

	// Err holds the error message if the step failed.
	Err string `json:"error,omitempty"`
}
