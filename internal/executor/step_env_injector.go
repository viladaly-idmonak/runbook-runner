package executor

import (
	"fmt"
	"strings"

	"github.com/user/runbook-runner/internal/config"
)

// StepEnvInjector merges per-step environment overrides into a base
// environment slice before the step command is executed.
type StepEnvInjector struct {
	cfg config.StepEnvConfig
}

// NewStepEnvInjector returns a StepEnvInjector backed by cfg.
func NewStepEnvInjector(cfg config.StepEnvConfig) *StepEnvInjector {
	return &StepEnvInjector{cfg: cfg}
}

// Inject returns a new env slice with any overrides for stepName merged
// on top of base. Overrides take precedence; base entries with the same
// key are replaced. When the injector is disabled or no overrides exist
// for stepName the original base slice is returned unchanged.
func (s *StepEnvInjector) Inject(stepName string, base []string) []string {
	if !s.cfg.Enabled {
		return base
	}
	overrides, ok := s.cfg.Overrides[stepName]
	if !ok || len(overrides) == 0 {
		return base
	}

	// Build a map from existing base so we can replace keys.
	result := make([]string, 0, len(base)+len(overrides))
	overrideKeys := overrideKeySet(overrides)

	for _, entry := range base {
		key := envKey(entry)
		if _, shadowed := overrideKeys[key]; !shadowed {
			result = append(result, entry)
		}
	}
	result = append(result, overrides...)
	return result
}

// Pairs returns the override pairs for a given step, or nil.
func (s *StepEnvInjector) Pairs(stepName string) []string {
	if !s.cfg.Enabled {
		return nil
	}
	return s.cfg.Overrides[stepName]
}

func overrideKeySet(pairs []string) map[string]struct{} {
	m := make(map[string]struct{}, len(pairs))
	for _, p := range pairs {
		m[envKey(p)] = struct{}{}
	}
	return m
}

func envKey(pair string) string {
	if idx := strings.IndexByte(pair, '='); idx >= 0 {
		return pair[:idx]
	}
	return fmt.Sprintf("__invalid_%s", pair)
}
