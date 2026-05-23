package config

import (
	"os"
	"strings"
	"testing"
)

func TestDefaultEnvConfig_Values(t *testing.T) {
	cfg := DefaultEnvConfig()
	if len(cfg.Passthrough) == 0 {
		t.Fatal("expected non-empty default passthrough list")
	}
	if cfg.Extra == nil {
		t.Fatal("expected Extra map to be initialised")
	}
	if cfg.AllowOverride {
		t.Fatal("expected AllowOverride to default to false")
	}
}

func TestValidateEnv_Valid(t *testing.T) {
	cfg := EnvConfig{
		Passthrough: []string{"HOME", "PATH"},
		Extra:       map[string]string{"FOO": "bar"},
	}
	if err := ValidateEnv(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateEnv_EmptyPassthroughName(t *testing.T) {
	cfg := EnvConfig{Passthrough: []string{"HOME", ""}}
	if err := ValidateEnv(cfg); err == nil {
		t.Fatal("expected error for empty passthrough name")
	}
}

func TestValidateEnv_InvalidPassthroughName(t *testing.T) {
	cfg := EnvConfig{Passthrough: []string{"MY VAR"}}
	if err := ValidateEnv(cfg); err == nil {
		t.Fatal("expected error for passthrough name with space")
	}
}

func TestValidateEnv_EmptyExtraKey(t *testing.T) {
	cfg := EnvConfig{Extra: map[string]string{"": "value"}}
	if err := ValidateEnv(cfg); err == nil {
		t.Fatal("expected error for empty extra key")
	}
}

func TestValidateEnv_ExtraKeyWithEquals(t *testing.T) {
	cfg := EnvConfig{Extra: map[string]string{"KEY=BAD": "val"}}
	if err := ValidateEnv(cfg); err == nil {
		t.Fatal("expected error for extra key containing '='")
	}
}

func TestResolve_PassthroughFromHost(t *testing.T) {
	t.Setenv("RUNBOOK_TEST_VAR", "hello")
	cfg := EnvConfig{
		Passthrough: []string{"RUNBOOK_TEST_VAR"},
		Extra:       map[string]string{},
	}
	env := cfg.Resolve()
	if !containsEntry(env, "RUNBOOK_TEST_VAR=hello") {
		t.Fatalf("expected RUNBOOK_TEST_VAR=hello in %v", env)
	}
}

func TestResolve_ExtraInjected(t *testing.T) {
	cfg := EnvConfig{
		Passthrough: []string{},
		Extra:       map[string]string{"STAGE": "test"},
	}
	env := cfg.Resolve()
	if !containsEntry(env, "STAGE=test") {
		t.Fatalf("expected STAGE=test in %v", env)
	}
}

func TestResolve_NoOverrideByDefault(t *testing.T) {
	os.Setenv("MYVAR", "original")
	defer os.Unsetenv("MYVAR")
	cfg := EnvConfig{
		Passthrough:   []string{"MYVAR"},
		Extra:         map[string]string{"MYVAR": "overridden"},
		AllowOverride: false,
	}
	env := cfg.Resolve()
	if !containsEntry(env, "MYVAR=original") {
		t.Fatalf("expected original value to be preserved, got %v", env)
	}
}

func TestResolve_AllowOverride(t *testing.T) {
	os.Setenv("MYVAR", "original")
	defer os.Unsetenv("MYVAR")
	cfg := EnvConfig{
		Passthrough:   []string{"MYVAR"},
		Extra:         map[string]string{"MYVAR": "overridden"},
		AllowOverride: true,
	}
	env := cfg.Resolve()
	if !containsEntry(env, "MYVAR=overridden") {
		t.Fatalf("expected overridden value, got %v", env)
	}
}

func containsEntry(env []string, entry string) bool {
	for _, e := range env {
		if strings.EqualFold(e, entry) || e == entry {
			return true
		}
	}
	return false
}
