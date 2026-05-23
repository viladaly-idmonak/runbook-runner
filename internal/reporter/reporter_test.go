package reporter_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/runbook-runner/internal/reporter"
)

func TestStepStatus_String(t *testing.T) {
	cases := []struct {
		status reporter.StepStatus
		want   string
	}{
		{reporter.StatusSuccess, "SUCCESS"},
		{reporter.StatusFailed, "FAILED"},
		{reporter.StatusSkipped, "SKIPPED"},
		{reporter.StatusRolledBack, "ROLLED_BACK"},
		{reporter.StatusPending, "PENDING"},
	}
	for _, tc := range cases {
		if got := tc.status.String(); got != tc.want {
			t.Errorf("StepStatus(%d).String() = %q, want %q", tc.status, got, tc.want)
		}
	}
}

func TestReporter_PrintSummary_ContainsTitle(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf)
	r.PrintSummary("Deploy App", []reporter.StepResult{})
	if !strings.Contains(buf.String(), "Deploy App") {
		t.Errorf("expected summary to contain runbook title, got:\n%s", buf.String())
	}
}

func TestReporter_PrintSummary_StepLines(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf)

	results := []reporter.StepResult{
		{Name: "Check disk", Status: reporter.StatusSuccess, Duration: 120 * time.Millisecond},
		{Name: "Run migration", Status: reporter.StatusFailed, Duration: 45 * time.Millisecond, Err: fmt.Errorf("exit code 1")},
		{Name: "Restart service", Status: reporter.StatusSkipped, Duration: 0},
	}

	r.PrintSummary("Maintenance", results)
	out := buf.String()

	for _, name := range []string{"Check disk", "Run migration", "Restart service"} {
		if !strings.Contains(out, name) {
			t.Errorf("expected output to contain step %q", name)
		}
	}
	if !strings.Contains(out, "FAILED") {
		t.Error("expected output to contain FAILED status")
	}
}

func TestReporter_PrintSummary_Tally(t *testing.T) {
	var buf bytes.Buffer
	r := reporter.New(&buf)

	results := []reporter.StepResult{
		{Name: "s1", Status: reporter.StatusSuccess},
		{Name: "s2", Status: reporter.StatusSuccess},
		{Name: "s3", Status: reporter.StatusFailed},
		{Name: "s4", Status: reporter.StatusRolledBack},
	}

	r.PrintSummary("Tally Test", results)
	out := buf.String()

	expected := []string{"Total: 4", "Success: 2", "Failed: 1", "RolledBack: 1"}
	for _, e := range expected {
		if !strings.Contains(out, e) {
			t.Errorf("expected tally to contain %q, got:\n%s", e, out)
		}
	}
}

func TestNew_DefaultsToStdout(t *testing.T) {
	// Ensure New(nil) does not panic.
	r := reporter.New(nil)
	if r == nil {
		t.Fatal("expected non-nil reporter")
	}
}
