package config

import "testing"

func TestLabelFilter_NoRules_AllowsAll(t *testing.T) {
	f := NewLabelFilter(DefaultLabelConfig())
	if !f.Allow(StepLabels{"env": "prod"}) {
		t.Error("expected step to be allowed when no rules are set")
	}
}

func TestLabelFilter_IncludeMatch_Allows(t *testing.T) {
	cfg := DefaultLabelConfig()
	cfg.Include = map[string]string{"env": "prod"}
	f := NewLabelFilter(cfg)
	if !f.Allow(StepLabels{"env": "prod", "team": "ops"}) {
		t.Error("expected step with matching include label to be allowed")
	}
}

func TestLabelFilter_IncludeMismatch_Blocks(t *testing.T) {
	cfg := DefaultLabelConfig()
	cfg.Include = map[string]string{"env": "prod"}
	f := NewLabelFilter(cfg)
	if f.Allow(StepLabels{"env": "staging"}) {
		t.Error("expected step with non-matching include label to be blocked")
	}
}

func TestLabelFilter_ExcludeMatch_Blocks(t *testing.T) {
	cfg := DefaultLabelConfig()
	cfg.Exclude = map[string]string{"skip": "true"}
	f := NewLabelFilter(cfg)
	if f.Allow(StepLabels{"skip": "true"}) {
		t.Error("expected step matching exclude label to be blocked")
	}
}

func TestLabelFilter_ExcludeMismatch_Allows(t *testing.T) {
	cfg := DefaultLabelConfig()
	cfg.Exclude = map[string]string{"skip": "true"}
	f := NewLabelFilter(cfg)
	if !f.Allow(StepLabels{"skip": "false"}) {
		t.Error("expected step with non-matching exclude label to be allowed")
	}
}

func TestLabelFilter_IncludeAndExclude_ExcludeTakesPrecedence(t *testing.T) {
	cfg := DefaultLabelConfig()
	cfg.Include = map[string]string{"env": "prod"}
	cfg.Exclude = map[string]string{"skip": "true"}
	f := NewLabelFilter(cfg)
	// matches include but also matches exclude → blocked
	if f.Allow(StepLabels{"env": "prod", "skip": "true"}) {
		t.Error("expected exclude to take precedence over include")
	}
	// matches include and does NOT match exclude → allowed
	if !f.Allow(StepLabels{"env": "prod", "skip": "false"}) {
		t.Error("expected step to be allowed when include matches and exclude does not")
	}
}

func TestLabelFilter_EmptyStepLabels_IncludeBlocks(t *testing.T) {
	cfg := DefaultLabelConfig()
	cfg.Include = map[string]string{"env": "prod"}
	f := NewLabelFilter(cfg)
	if f.Allow(StepLabels{}) {
		t.Error("expected step with no labels to be blocked when include is set")
	}
}
