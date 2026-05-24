package executor_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/runbook-runner/internal/config"
	"github.com/your-org/runbook-runner/internal/executor"
)

func makePrompter(cfg config.PromptConfig, input string) (*executor.Prompter, *bytes.Buffer) {
	out := &bytes.Buffer{}
	p := executor.NewPrompter(cfg, strings.NewReader(input), out)
	return p, out
}

func TestPrompter_DisabledAlwaysTrue(t *testing.T) {
	cfg := config.DefaultPromptConfig() // Enabled=false
	p, _ := makePrompter(cfg, "")
	ok, err := p.ConfirmStep("deploy")
	if err != nil || !ok {
		t.Errorf("expected (true, nil), got (%v, %v)", ok, err)
	}
}

func TestPrompter_NonInteractiveAlwaysTrue(t *testing.T) {
	cfg := config.PromptConfig{Enabled: false, NonInteractive: true}
	p, _ := makePrompter(cfg, "")
	ok, err := p.ConfirmStep("deploy")
	if err != nil || !ok {
		t.Errorf("expected (true, nil), got (%v, %v)", ok, err)
	}
}

func TestPrompter_ConfirmStep_Yes(t *testing.T) {
	cfg := config.PromptConfig{Enabled: true}
	p, _ := makePrompter(cfg, "y\n")
	ok, err := p.ConfirmStep("migrate")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("expected true for 'y' input")
	}
}

func TestPrompter_ConfirmStep_No(t *testing.T) {
	cfg := config.PromptConfig{Enabled: true}
	p, _ := makePrompter(cfg, "n\n")
	ok, err := p.ConfirmStep("migrate")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected false for 'n' input")
	}
}

func TestPrompter_ConfirmRollback_OnFailureDisabled(t *testing.T) {
	cfg := config.PromptConfig{OnFailure: false}
	p, _ := makePrompter(cfg, "")
	ok, err := p.ConfirmRollback("migrate")
	if err != nil || !ok {
		t.Errorf("expected (true, nil), got (%v, %v)", ok, err)
	}
}

func TestPrompter_ConfirmRollback_Yes(t *testing.T) {
	cfg := config.PromptConfig{OnFailure: true}
	p, _ := makePrompter(cfg, "yes\n")
	ok, err := p.ConfirmRollback("migrate")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("expected true for 'yes' input")
	}
}

func TestPrompter_EOFReturnsFalse(t *testing.T) {
	cfg := config.PromptConfig{Enabled: true}
	p, _ := makePrompter(cfg, "") // empty reader → EOF
	ok, err := p.ConfirmStep("cleanup")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected false on EOF")
	}
}
