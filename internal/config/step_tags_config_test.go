package config

import "testing"

func TestDefaultStepTagsConfig_Values(t *testing.T) {
	c := DefaultStepTagsConfig()
	if c.Enabled {
		t.Error("expected Enabled to be false")
	}
	if c.StepTags == nil {
		t.Error("expected StepTags to be non-nil map")
	}
	if len(c.StepTags) != 0 {
		t.Errorf("expected empty StepTags, got %d entries", len(c.StepTags))
	}
}

func TestValidateStepTags_DisabledSkipsValidation(t *testing.T) {
	c := StepTagsConfig{
		Enabled:  false,
		StepTags: map[string][]string{"": {"tag"}},
	}
	if err := ValidateStepTags(c); err != nil {
		t.Errorf("expected no error when disabled, got %v", err)
	}
}

func TestValidateStepTags_Valid(t *testing.T) {
	c := StepTagsConfig{
		Enabled:  true,
		StepTags: map[string][]string{"deploy": {"prod", "infra"}},
	}
	if err := ValidateStepTags(c); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateStepTags_EmptyStepName(t *testing.T) {
	c := StepTagsConfig{
		Enabled:  true,
		StepTags: map[string][]string{"": {"tag"}},
	}
	if err := ValidateStepTags(c); err == nil {
		t.Error("expected error for empty step name")
	}
}

func TestValidateStepTags_EmptyTag(t *testing.T) {
	c := StepTagsConfig{
		Enabled:  true,
		StepTags: map[string][]string{"build": {"ok", ""}},
	}
	if err := ValidateStepTags(c); err == nil {
		t.Error("expected error for empty tag value")
	}
}

func TestTagsForStep_Disabled(t *testing.T) {
	c := StepTagsConfig{
		Enabled:  false,
		StepTags: map[string][]string{"build": {"fast"}},
	}
	if tags := TagsForStep(c, "build"); tags != nil {
		t.Errorf("expected nil when disabled, got %v", tags)
	}
}

func TestTagsForStep_ReturnsCorrectTags(t *testing.T) {
	c := StepTagsConfig{
		Enabled:  true,
		StepTags: map[string][]string{"deploy": {"prod", "critical"}},
	}
	tags := TagsForStep(c, "deploy")
	if len(tags) != 2 || tags[0] != "prod" || tags[1] != "critical" {
		t.Errorf("unexpected tags: %v", tags)
	}
}

func TestTagsForStep_MissingStep(t *testing.T) {
	c := StepTagsConfig{
		Enabled:  true,
		StepTags: map[string][]string{},
	}
	if tags := TagsForStep(c, "missing"); tags != nil {
		t.Errorf("expected nil for unknown step, got %v", tags)
	}
}
