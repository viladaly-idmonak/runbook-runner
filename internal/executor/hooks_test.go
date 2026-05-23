package executor

import (
	"runtime"
	"testing"
)

func TestHookRunner_NoMatchingHooks(t *testing.T) {
	hr := NewHookRunner("/bin/sh", []Hook{}, false)
	if err := hr.Run(HookPreStep, 0, "step-one"); err != nil {
		t.Fatalf("expected no error with empty hooks, got %v", err)
	}
}

func TestHookRunner_PreStepHookRuns(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping shell test on windows")
	}
	hooks := []Hook{
		{Type: HookPreStep, Command: "exit 0"},
	}
	hr := NewHookRunner("/bin/sh", hooks, false)
	if err := hr.Run(HookPreStep, 1, "deploy"); err != nil {
		t.Fatalf("expected success, got %v", err)
	}
}

func TestHookRunner_FailingHookReturnsError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping shell test on windows")
	}
	hooks := []Hook{
		{Type: HookOnError, Command: "exit 1"},
	}
	hr := NewHookRunner("/bin/sh", hooks, false)
	err := hr.Run(HookOnError, 2, "migrate")
	if err == nil {
		t.Fatal("expected error from failing hook")
	}
}

func TestHookRunner_WrongTypeSkipped(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping shell test on windows")
	}
	hooks := []Hook{
		{Type: HookPostStep, Command: "exit 1"},
	}
	hr := NewHookRunner("/bin/sh", hooks, false)
	// pre-step hooks should not trigger the post-step failing hook
	if err := hr.Run(HookPreStep, 0, "init"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestHookRunner_MultipleHooksSameType(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping shell test on windows")
	}
	hooks := []Hook{
		{Type: HookPreStep, Command: "exit 0"},
		{Type: HookPreStep, Command: "exit 0"},
	}
	hr := NewHookRunner("/bin/sh", hooks, false)
	if err := hr.Run(HookPreStep, 0, "start"); err != nil {
		t.Fatalf("expected success for multiple passing hooks, got %v", err)
	}
}
