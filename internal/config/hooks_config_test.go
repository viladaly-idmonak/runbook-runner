package config

import (
	"errors"
	"testing"
)

func TestValidateHooks_Empty(t *testing.T) {
	if err := ValidateHooks(nil); err != nil {
		t.Fatalf("expected no error for empty hooks, got %v", err)
	}
}

func TestValidateHooks_ValidTypes(t *testing.T) {
	hooks := []HookConfig{
		{Type: "pre-step", Command: "echo pre"},
		{Type: "post-step", Command: "echo post"},
		{Type: "on-error", Command: "echo err"},
	}
	if err := ValidateHooks(hooks); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateHooks_InvalidType(t *testing.T) {
	hooks := []HookConfig{
		{Type: "before-all", Command: "echo x"},
	}
	err := ValidateHooks(hooks)
	if err == nil {
		t.Fatal("expected error for unknown hook type")
	}
	if !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("expected ErrInvalidConfig, got %v", err)
	}
}

func TestValidateHooks_EmptyCommand(t *testing.T) {
	hooks := []HookConfig{
		{Type: "pre-step", Command: ""},
	}
	err := ValidateHooks(hooks)
	if err == nil {
		t.Fatal("expected error for empty command")
	}
	if !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("expected ErrInvalidConfig, got %v", err)
	}
}

func TestValidateHooks_SecondEntryInvalid(t *testing.T) {
	hooks := []HookConfig{
		{Type: "pre-step", Command: "echo ok"},
		{Type: "unknown", Command: "echo bad"},
	}
	if err := ValidateHooks(hooks); err == nil {
		t.Fatal("expected error for second invalid hook")
	}
}
