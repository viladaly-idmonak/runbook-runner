package reporter

import (
	"encoding/json"
	"io"
	"time"
)

// JSONReport is the top-level structure emitted when output format is "json".
type JSONReport struct {
	Title     string      `json:"title"`
	StartedAt time.Time   `json:"started_at"`
	Steps     []JSONStep  `json:"steps"`
	Passed    int         `json:"passed"`
	Failed    int         `json:"failed"`
	Skipped   int         `json:"skipped"`
}

// JSONStep captures the outcome of a single runbook step.
type JSONStep struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Output string `json:"output,omitempty"`
	Error  string `json:"error,omitempty"`
}

// JSONReporter writes a single JSON object to w when Flush is called.
type JSONReporter struct {
	w       io.Writer
	report  JSONReport
}

// NewJSONReporter creates a JSONReporter that writes to w.
func NewJSONReporter(title string, w io.Writer) *JSONReporter {
	return &JSONReporter{
		w: w,
		report: JSONReport{
			Title:     title,
			StartedAt: time.Now(),
		},
	}
}

// Record appends a step result to the in-memory report.
func (r *JSONReporter) Record(name, status, output, errMsg string) {
	r.report.Steps = append(r.report.Steps, JSONStep{
		Name:   name,
		Status: status,
		Output: output,
		Error:  errMsg,
	})
	switch status {
	case "ok":
		r.report.Passed++
	case "fail":
		r.report.Failed++
	case "skip":
		r.report.Skipped++
	}
}

// Flush encodes the report as indented JSON and writes it to the writer.
func (r *JSONReporter) Flush() error {
	enc := json.NewEncoder(r.w)
	enc.SetIndent("", "  ")
	return enc.Encode(r.report)
}
