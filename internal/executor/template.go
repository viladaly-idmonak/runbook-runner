package executor

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/user/runbook-runner/internal/config"
)

// TemplateRenderer applies variable substitution to step command strings.
type TemplateRenderer struct {
	cfg  config.TemplateConfig
	vars map[string]string
}

// NewTemplateRenderer creates a renderer from the given config.
// If cfg.Enabled is false, Render is a no-op passthrough.
func NewTemplateRenderer(cfg config.TemplateConfig) *TemplateRenderer {
	vars := make(map[string]string, len(cfg.Vars))
	for k, v := range cfg.Vars {
		vars[k] = v
	}
	return &TemplateRenderer{cfg: cfg, vars: vars}
}

// SetVar adds or overrides a single variable at runtime.
func (r *TemplateRenderer) SetVar(key, value string) {
	r.vars[key] = value
}

// Render substitutes template variables in src and returns the result.
// When cfg.Enabled is false the original string is returned unchanged.
func (r *TemplateRenderer) Render(src string) (string, error) {
	if !r.cfg.Enabled {
		return src, nil
	}

	option := "missingkey=zero"
	if r.cfg.StrictMode {
		option = "missingkey=error"
	}

	tmpl, err := template.New("").
		Delims(r.cfg.LeftDelim, r.cfg.RightDelim).
		Option(option).
		Parse(src)
	if err != nil {
		return "", fmt.Errorf("template parse: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, r.vars); err != nil {
		return "", fmt.Errorf("template execute: %w", err)
	}
	return strings.TrimRight(buf.String(), "\x00"), nil
}
