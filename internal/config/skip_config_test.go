package config

import (
	"testing"
)

func TestDefaultSkipConfig_Values(t *testing.T) {
	c := DefaultSkipConfig()
	if len(c.StepNames) != 0 {
		t.Errorf("expected empty StepNames, got %v", c.StepNames)
	}
	if len(c.StepIndices) != 0 {
		t.Errorf("expected empty StepIndices, got %v", c.StepIndices)
	}
	if c.OnlyWhen != "" {
		t.Errorf("expected empty OnlyWhen, got %q", c.OnlyWhen)
	}
}

func TestValidateSkip_Valid(t *testing.T) {
	c := SkipConfig{
		StepNames:   []string{"setup", "teardown"},
		StepIndices: []int{0, 2, 4},
	}
	if err := ValidateSkip(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateSkip_EmptyStepName(t *testing.T) {
	c := SkipConfig{
		StepNames: []string{"ok", ""},
	}
	if err := ValidateSkip(c); err == nil {
		t.Fatal("expected error for empty step name, got nil")
	}
}

func TestValidateSkip_DuplicateStepName(t *testing.T) {
	c := SkipConfig{
		StepNames: []string{"deploy", "deploy"},
	}
	if err := ValidateSkip(c); err == nil {
		t.Fatal("expected error for duplicate step name, got nil")
	}
}

func TestValidateSkip_NegativeIndex(t *testing.T) {
	c := SkipConfig{
		StepIndices: []int{1, -1},
	}
	if err := ValidateSkip(c); err == nil {
		t.Fatal("expected error for negative index, got nil")
	}
}

func TestValidateSkip_DuplicateIndex(t *testing.T) {
	c := SkipConfig{
		StepIndices: []int{3, 3},
	}
	if err := ValidateSkip(c); err == nil {
		t.Fatal("expected error for duplicate index, got nil")
	}
}

func TestShouldSkip_ByName(t *testing.T) {
	c := SkipConfig{StepNames: []string{"migrate"}}
	if !c.ShouldSkip(0, "migrate") {
		t.Error("expected ShouldSkip=true for matching name")
	}
	if c.ShouldSkip(0, "rollback") {
		t.Error("expected ShouldSkip=false for non-matching name")
	}
}

func TestShouldSkip_ByIndex(t *testing.T) {
	c := SkipConfig{StepIndices: []int{2, 5}}
	if !c.ShouldSkip(2, "anything") {
		t.Error("expected ShouldSkip=true for matching index")
	}
	if c.ShouldSkip(3, "anything") {
		t.Error("expected ShouldSkip=false for non-matching index")
	}
}

func TestShouldSkip_EmptyConfig(t *testing.T) {
	c := DefaultSkipConfig()
	if c.ShouldSkip(0, "setup") {
		t.Error("expected ShouldSkip=false with empty config")
	}
}
