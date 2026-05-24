package executor_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/runbook-runner/internal/config"
	"github.com/your-org/runbook-runner/internal/executor"
)

// TestPrompter_Integration_PrintsQuestion verifies the question is written to out.
func TestPrompter_Integration_PrintsQuestion(t *testing.T) {
	cfg := config.PromptConfig{Enabled: true}
	out := &bytes.Buffer{}
	p := executor.NewPrompter(cfg, strings.NewReader("y\n"), out)

	_, err := p.ConfirmStep("restart-service")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "restart-service") {
		t.Errorf("expected question to contain step name, got: %q", out.String())
	}
}

// TestPrompter_Integration_MultipleSteps simulates confirming two steps in sequence.
func TestPrompter_Integration_MultipleSteps(t *testing.T) {
	cfg := config.PromptConfig{Enabled: true}
	out := &bytes.Buffer{}
	p := executor.NewPrompter(cfg, strings.NewReader("y\nn\n"), out)

	ok1, err := p.ConfirmStep("step-one")
	if err != nil || !ok1 {
		t.Errorf("step-one: expected (true, nil), got (%v, %v)", ok1, err)
	}

	ok2, err := p.ConfirmStep("step-two")
	if err != nil || ok2 {
		t.Errorf("step-two: expected (false, nil), got (%v, %v)", ok2, err)
	}
}

// TestPrompter_Integration_RollbackFlow simulates a rollback confirmation.
func TestPrompter_Integration_RollbackFlow(t *testing.T) {
	cfg := config.PromptConfig{Enabled: true, OnFailure: true}
	out := &bytes.Buffer{}
	p := executor.NewPrompter(cfg, strings.NewReader("y\nyes\n"), out)

	okStep, _ := p.ConfirmStep("deploy")
	if !okStep {
		t.Error("expected step confirmation")
	}

	okRoll, _ := p.ConfirmRollback("deploy")
	if !okRoll {
		t.Error("expected rollback confirmation")
	}

	if !strings.Contains(out.String(), "rollback") {
		t.Errorf("expected rollback prompt in output, got: %q", out.String())
	}
}
