package executor

import (
	"bytes"
	"strings"

	"github.com/your-org/runbook-runner/internal/config"
)

// StepOutputCapture captures stdout from a step execution and stores it
// under a named variable for use in subsequent template rendering.
type StepOutputCapture struct {
	cfg    config.StepOutputConfig
	values map[string]string
}

// NewStepOutputCapture creates a new StepOutputCapture using the provided config.
func NewStepOutputCapture(cfg config.StepOutputConfig) *StepOutputCapture {
	return &StepOutputCapture{
		cfg:    cfg,
		values: make(map[string]string),
	}
}

// Record stores the raw output for the given step name, applying
// MaxBytes truncation and optional whitespace trimming.
func (c *StepOutputCapture) Record(stepName string, raw []byte) {
	if !c.cfg.Enabled {
		return
	}
	data := raw
	if c.cfg.MaxBytes > 0 && len(data) > c.cfg.MaxBytes {
		data = data[:c.cfg.MaxBytes]
	}
	value := string(data)
	if c.cfg.TrimSpace {
		value = strings.TrimSpace(value)
	}
	varName := c.varName(stepName)
	c.values[varName] = value
}

// Get returns the captured output variable value for a given variable name.
func (c *StepOutputCapture) Get(varName string) (string, bool) {
	v, ok := c.values[varName]
	return v, ok
}

// All returns a copy of all captured variable name -> value pairs.
func (c *StepOutputCapture) All() map[string]string {
	out := make(map[string]string, len(c.values))
	for k, v := range c.values {
		out[k] = v
	}
	return out
}

// varName resolves the variable name for a step, using StepOverrides if present,
// otherwise deriving a default from the step name.
func (c *StepOutputCapture) varName(stepName string) string {
	if override, ok := c.cfg.StepOverrides[stepName]; ok {
		return override
	}
	// default: STEP_<UPPER_SNAKE> e.g. "run deploy" -> "STEP_RUN_DEPLOY"
	replacer := strings.NewReplacer(" ", "_", "-", "_")
	return "STEP_" + strings.ToUpper(replacer.Replace(stepName))
}

// LimitedBuffer is a bytes.Buffer that stops writing after MaxBytes.
type LimitedBuffer struct {
	buf      bytes.Buffer
	maxBytes int
	written  int
}

// Write appends bytes up to the configured limit.
func (lb *LimitedBuffer) Write(p []byte) (int, error) {
	if lb.maxBytes > 0 {
		remaining := lb.maxBytes - lb.written
		if remaining <= 0 {
			return len(p), nil
		}
		if len(p) > remaining {
			p = p[:remaining]
		}
	}
	n, err := lb.buf.Write(p)
	lb.written += n
	return n, err
}

// Bytes returns the buffered bytes.
func (lb *LimitedBuffer) Bytes() []byte {
	return lb.buf.Bytes()
}
