package reporter

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// StepStatus represents the outcome of a single runbook step.
type StepStatus int

const (
	StatusPending StepStatus = iota
	StatusSuccess
	StatusFailed
	StatusSkipped
	StatusRolledBack
)

func (s StepStatus) String() string {
	switch s {
	case StatusSuccess:
		return "SUCCESS"
	case StatusFailed:
		return "FAILED"
	case StatusSkipped:
		return "SKIPPED"
	case StatusRolledBack:
		return "ROLLED_BACK"
	default:
		return "PENDING"
	}
}

// StepResult holds the result of executing a single step.
type StepResult struct {
	Name     string
	Status   StepStatus
	Output   string
	Err      error
	Duration time.Duration
}

// Reporter writes structured execution summaries.
type Reporter struct {
	out io.Writer
}

// New creates a Reporter writing to the given writer.
// If w is nil, os.Stdout is used.
func New(w io.Writer) *Reporter {
	if w == nil {
		w = os.Stdout
	}
	return &Reporter{out: w}
}

// PrintSummary writes a formatted summary of all step results.
func (r *Reporter) PrintSummary(runbookTitle string, results []StepResult) {
	fmt.Fprintf(r.out, "\n%s\n", strings.Repeat("=", 50))
	fmt.Fprintf(r.out, "Runbook: %s\n", runbookTitle)
	fmt.Fprintf(r.out, "%s\n\n", strings.Repeat("=", 50))

	for i, res := range results {
		fmt.Fprintf(r.out, "[%d] %-30s %s (%s)\n",
			i+1, res.Name, res.Status, res.Duration.Round(time.Millisecond))
		if res.Err != nil {
			fmt.Fprintf(r.out, "    error: %v\n", res.Err)
		}
	}

	succeeded, failed, skipped, rolledBack := tally(results)
	fmt.Fprintf(r.out, "\n%s\n", strings.Repeat("-", 50))
	fmt.Fprintf(r.out, "Total: %d  Success: %d  Failed: %d  Skipped: %d  RolledBack: %d\n",
		len(results), succeeded, failed, skipped, rolledBack)
}

func tally(results []StepResult) (succeeded, failed, skipped, rolledBack int) {
	for _, r := range results {
		switch r.Status {
		case StatusSuccess:
			succeeded++
		case StatusFailed:
			failed++
		case StatusSkipped:
			skipped++
		case StatusRolledBack:
			rolledBack++
		}
	}
	return
}
