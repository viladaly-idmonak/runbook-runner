package executor

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/runbook-runner/internal/config"
)

// EnvResolver builds the environment slice for a step's execution.
// It merges the base process environment with extra key/value pairs from
// config and selectively passes through variables listed in Passthrough.
type EnvResolver struct {
	cfg config.EnvConfig
}

// NewEnvResolver creates an EnvResolver backed by the given EnvConfig.
func NewEnvResolver(cfg config.EnvConfig) *EnvResolver {
	return &EnvResolver{cfg: cfg}
}

// Resolve returns the environment that should be used when running a step.
// If InheritAll is true the full os.Environ is included; otherwise only
// variables named in Passthrough are forwarded from the host environment.
// Extra key=value pairs from config are always appended and take precedence.
func (r *EnvResolver) Resolve() []string {
	var env []string

	if r.cfg.InheritAll {
		env = append(env, os.Environ()...)
	} else {
		for _, name := range r.cfg.Passthrough {
			if val, ok := os.LookupEnv(name); ok {
				env = append(env, fmt.Sprintf("%s=%s", name, val))
			}
		}
	}

	for k, v := range r.cfg.Extra {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	return env
}

// ToMap converts a slice of "KEY=VALUE" strings into a map for easy lookup.
func ToMap(env []string) map[string]string {
	m := make(map[string]string, len(env))
	for _, e := range env {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			m[parts[0]] = parts[1]
		}
	}
	return m
}
