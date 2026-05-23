package executor

import (
	"fmt"
	"io"
	"os"
	"time"
)

// StepLogger records per-step execution details including timing and output.
type StepLogger struct {
	out    io.Writer
	verbose bool
}

// StepLog holds the result of a single step execution.
type StepLog struct {
	StepName  string
	Command   string
	Output    string
	Err       error
	StartedAt time.Time
	Duration  time.Duration
	Skipped   bool
}

// NewStepLogger creates a StepLogger writing to w. If w is nil, os.Stdout is used.
func NewStepLogger(w io.Writer, verbose bool) *StepLogger {
	if w == nil {
		w = os.Stdout
	}
	return &StepLogger{out: w, verbose: verbose}
}

// Log writes a formatted entry for the given StepLog.
func (sl *StepLogger) Log(entry StepLog) {
	status := "OK"
	if entry.Skipped {
		status = "SKIP"
	} else if entry.Err != nil {
		status = "FAIL"
	}

	fmt.Fprintf(sl.out, "[%s] %-40s %s (%s)\n",
		status,
		entry.StepName,
		entry.StartedAt.Format(time.RFC3339),
		entry.Duration.Round(time.Millisecond),
	)

	if sl.verbose && entry.Output != "" {
		fmt.Fprintf(sl.out, "       output: %s\n", entry.Output)
	}

	if entry.Err != nil {
		fmt.Fprintf(sl.out, "       error:  %v\n", entry.Err)
	}
}
