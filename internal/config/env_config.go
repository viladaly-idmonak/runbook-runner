package config

import (
	"fmt"
	"os"
	"strings"
)

// EnvConfig holds environment variable injection settings for runbook execution.
type EnvConfig struct {
	// Passthrough lists environment variable names to forward from the host.
	Passthrough []string `yaml:"passthrough"`
	// Extra holds additional key=value pairs to inject into each step's environment.
	Extra map[string]string `yaml:"extra"`
	// AllowOverride controls whether Extra values can overwrite host env vars.
	AllowOverride bool `yaml:"allow_override"`
}

// DefaultEnvConfig returns a safe default EnvConfig.
func DefaultEnvConfig() EnvConfig {
	return EnvConfig{
		Passthrough:   []string{"PATH", "HOME", "USER"},
		Extra:         map[string]string{},
		AllowOverride: false,
	}
}

// ValidateEnv checks an EnvConfig for invalid entries.
func ValidateEnv(e EnvConfig) error {
	for _, name := range e.Passthrough {
		if strings.TrimSpace(name) == "" {
			return fmt.Errorf("env: passthrough contains an empty variable name")
		}
		if strings.ContainsAny(name, "= ") {
			return fmt.Errorf("env: passthrough variable name %q must not contain '=' or spaces", name)
		}
	}
	for k, v := range e.Extra {
		if strings.TrimSpace(k) == "" {
			return fmt.Errorf("env: extra contains an empty key")
		}
		if strings.ContainsAny(k, "= ") {
			return fmt.Errorf("env: extra key %q must not contain '=' or spaces", k)
		}
		_ = v // values are unrestricted
	}
	return nil
}

// Resolve builds the final environment slice to pass to a subprocess.
// It starts with passthrough vars from the host, then merges Extra according
// to AllowOverride.
func (e EnvConfig) Resolve() []string {
	env := make(map[string]string, len(e.Passthrough)+len(e.Extra))

	for _, name := range e.Passthrough {
		if val, ok := os.LookupEnv(name); ok {
			env[name] = val
		}
	}

	for k, v := range e.Extra {
		if _, exists := env[k]; exists && !e.AllowOverride {
			continue
		}
		env[k] = v
	}

	out := make([]string, 0, len(env))
	for k, v := range env {
		out = append(out, k+"="+v)
	}
	return out
}
