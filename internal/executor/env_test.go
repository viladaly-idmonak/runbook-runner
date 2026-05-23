package executor

import (
	"os"
	"testing"

	"github.com/user/runbook-runner/internal/config"
)

func TestEnvResolver_InheritAll(t *testing.T) {
	os.Setenv("_RR_TEST_VAR", "hello")
	t.Cleanup(func() { os.Unsetenv("_RR_TEST_VAR") })

	cfg := config.DefaultEnvConfig()
	cfg.InheritAll = true

	r := NewEnvResolver(cfg)
	env := ToMap(r.Resolve())

	if env["_RR_TEST_VAR"] != "hello" {
		t.Errorf("expected _RR_TEST_VAR=hello, got %q", env["_RR_TEST_VAR"])
	}
}

func TestEnvResolver_Passthrough(t *testing.T) {
	os.Setenv("_RR_PASS", "passed")
	os.Setenv("_RR_SKIP", "skipped")
	t.Cleanup(func() {
		os.Unsetenv("_RR_PASS")
		os.Unsetenv("_RR_SKIP")
	})

	cfg := config.DefaultEnvConfig()
	cfg.InheritAll = false
	cfg.Passthrough = []string{"_RR_PASS"}

	r := NewEnvResolver(cfg)
	env := ToMap(r.Resolve())

	if env["_RR_PASS"] != "passed" {
		t.Errorf("expected _RR_PASS=passed, got %q", env["_RR_PASS"])
	}
	if _, ok := env["_RR_SKIP"]; ok {
		t.Error("_RR_SKIP should not be in resolved env")
	}
}

func TestEnvResolver_ExtraOverrides(t *testing.T) {
	os.Setenv("_RR_EXTRA", "original")
	t.Cleanup(func() { os.Unsetenv("_RR_EXTRA") })

	cfg := config.DefaultEnvConfig()
	cfg.InheritAll = true
	cfg.Extra = map[string]string{"_RR_EXTRA": "overridden"}

	r := NewEnvResolver(cfg)
	env := ToMap(r.Resolve())

	// Extra entries are appended last; ToMap keeps the last value for a key.
	if env["_RR_EXTRA"] != "overridden" {
		t.Errorf("expected _RR_EXTRA=overridden, got %q", env["_RR_EXTRA"])
	}
}

func TestEnvResolver_NoInheritNoPassthrough(t *testing.T) {
	cfg := config.DefaultEnvConfig()
	cfg.InheritAll = false
	cfg.Passthrough = nil
	cfg.Extra = map[string]string{"ONLY": "this"}

	r := NewEnvResolver(cfg)
	env := ToMap(r.Resolve())

	if len(env) != 1 || env["ONLY"] != "this" {
		t.Errorf("unexpected env: %v", env)
	}
}

func TestToMap_EmptySlice(t *testing.T) {
	m := ToMap([]string{})
	if len(m) != 0 {
		t.Errorf("expected empty map, got %v", m)
	}
}

func TestToMap_MalformedEntry(t *testing.T) {
	m := ToMap([]string{"NOEQUALSIGN"})
	if len(m) != 0 {
		t.Errorf("expected malformed entry to be skipped, got %v", m)
	}
}
