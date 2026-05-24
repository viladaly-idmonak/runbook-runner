package executor

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/your-org/runbook-runner/internal/config"
)

// Prompter asks the user for confirmation before proceeding.
type Prompter struct {
	cfg    config.PromptConfig
	in     io.Reader
	out    io.Writer
}

// NewPrompter creates a Prompter backed by the given reader/writer.
func NewPrompter(cfg config.PromptConfig, in io.Reader, out io.Writer) *Prompter {
	return &Prompter{cfg: cfg, in: in, out: out}
}

// ConfirmStep asks the user whether to proceed with the named step.
// Returns true when the user confirms or when prompts are disabled / non-interactive.
func (p *Prompter) ConfirmStep(stepName string) (bool, error) {
	if !p.cfg.Enabled || p.cfg.NonInteractive {
		return true, nil
	}
	return p.ask(fmt.Sprintf("Run step %q? [y/N]: ", stepName))
}

// ConfirmRollback asks the user whether to run a rollback command.
// Returns true when the user confirms or when on_failure prompts are disabled.
func (p *Prompter) ConfirmRollback(stepName string) (bool, error) {
	if !p.cfg.OnFailure || p.cfg.NonInteractive {
		return true, nil
	}
	return p.ask(fmt.Sprintf("Run rollback for step %q? [y/N]: ", stepName))
}

func (p *Prompter) ask(question string) (bool, error) {
	fmt.Fprint(p.out, question)
	scanner := bufio.NewScanner(p.in)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return false, fmt.Errorf("prompt: read error: %w", err)
		}
		// EOF — treat as no
		return false, nil
	}
	answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
	return answer == "y" || answer == "yes", nil
}
