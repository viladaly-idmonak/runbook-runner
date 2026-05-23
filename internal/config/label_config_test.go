package config

import (
	"testing"
)

func TestDefaultLabelConfig_Values(t *testing.T) {
	cfg := DefaultLabelConfig()
	if len(cfg.IncludeLabels) != 0 {
		t.Errorf("expected empty IncludeLabels, got %v", cfg.IncludeLabels)
	}
	if len(cfg.ExcludeLabels) != 0 {
		t.Errorf("expected empty ExcludeLabels, got %v", cfg.ExcludeLabels)
	}
}

func TestValidateLabels_Valid(t *testing.T) {
	cfg := LabelConfig{
		IncludeLabels: []string{"smoke", "fast"},
		ExcludeLabels: []string{"slow"},
	}
	if err := ValidateLabels(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateLabels_EmptyIncludeEntry(t *testing.T) {
	cfg := LabelConfig{IncludeLabels: []string{"ok", ""}}
	if err := ValidateLabels(cfg); err == nil {
		t.Fatal("expected error for empty include label")
	}
}

func TestValidateLabels_EmptyExcludeEntry(t *testing.T) {
	cfg := LabelConfig{ExcludeLabels: []string{""}}
	if err := ValidateLabels(cfg); err == nil {
		t.Fatal("expected error for empty exclude label")
	}
}

func TestValidateLabels_ConflictingLabel(t *testing.T) {
	cfg := LabelConfig{
		IncludeLabels: []string{"smoke"},
		ExcludeLabels: []string{"smoke"},
	}
	if err := ValidateLabels(cfg); err == nil {
		t.Fatal("expected error for label in both include and exclude")
	}
}

func TestMatchesFilter_NoFilter(t *testing.T) {
	cfg := DefaultLabelConfig()
	if !MatchesFilter(cfg, []string{"anything"}) {
		t.Error("expected match when no filter is set")
	}
}

func TestMatchesFilter_IncludeHit(t *testing.T) {
	cfg := LabelConfig{IncludeLabels: []string{"smoke"}}
	if !MatchesFilter(cfg, []string{"smoke", "integration"}) {
		t.Error("expected match on included label")
	}
}

func TestMatchesFilter_IncludeMiss(t *testing.T) {
	cfg := LabelConfig{IncludeLabels: []string{"smoke"}}
	if MatchesFilter(cfg, []string{"integration"}) {
		t.Error("expected no match when step lacks included label")
	}
}

func TestMatchesFilter_ExcludeTakesPriority(t *testing.T) {
	cfg := LabelConfig{
		IncludeLabels: []string{"smoke"},
		ExcludeLabels: []string{"slow"},
	}
	// Step has both an included and an excluded label — exclusion wins.
	if MatchesFilter(cfg, []string{"smoke", "slow"}) {
		t.Error("expected exclusion to take priority over inclusion")
	}
}

func TestMatchesFilter_NoLabelsWithIncludeFilter(t *testing.T) {
	cfg := LabelConfig{IncludeLabels: []string{"smoke"}}
	if MatchesFilter(cfg, []string{}) {
		t.Error("expected no match for unlabelled step when include filter is set")
	}
}
