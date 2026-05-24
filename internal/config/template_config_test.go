package config

import (
	"testing"
)

func TestDefaultTemplateConfig_Values(t *testing.T) {
	cfg := DefaultTemplateConfig()
	if cfg.Enabled {
		t.Error("expected Enabled=false")
	}
	if cfg.LeftDelim != "{{" {
		t.Errorf("expected LeftDelim='{{', got %q", cfg.LeftDelim)
	}
	if cfg.RightDelim != "}}" {
		t.Errorf("expected RightDelim='}}', got %q", cfg.RightDelim)
	}
	if cfg.Vars == nil {
		t.Error("expected Vars to be initialised")
	}
	if cfg.StrictMode {
		t.Error("expected StrictMode=false")
	}
}

func TestValidateTemplate_DisabledSkipsValidation(t *testing.T) {
	cfg := TemplateConfig{Enabled: false, LeftDelim: "", RightDelim: ""}
	if err := ValidateTemplate(cfg); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestValidateTemplate_Valid(t *testing.T) {
	cfg := TemplateConfig{
		Enabled:    true,
		LeftDelim:  "{{",
		RightDelim: "}}",
		Vars:       map[string]string{"env": "prod"},
	}
	if err := ValidateTemplate(cfg); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateTemplate_EmptyLeftDelim(t *testing.T) {
	cfg := TemplateConfig{Enabled: true, LeftDelim: "", RightDelim: "}}"}
	if err := ValidateTemplate(cfg); err == nil {
		t.Error("expected error for empty left_delim")
	}
}

func TestValidateTemplate_EmptyRightDelim(t *testing.T) {
	cfg := TemplateConfig{Enabled: true, LeftDelim: "{{", RightDelim: ""}
	if err := ValidateTemplate(cfg); err == nil {
		t.Error("expected error for empty right_delim")
	}
}

func TestValidateTemplate_EqualDelims(t *testing.T) {
	cfg := TemplateConfig{Enabled: true, LeftDelim: "!!", RightDelim: "!!"}
	if err := ValidateTemplate(cfg); err == nil {
		t.Error("expected error when left_delim equals right_delim")
	}
}

func TestValidateTemplate_DelimWithWhitespace(t *testing.T) {
	cfg := TemplateConfig{Enabled: true, LeftDelim: "{ {", RightDelim: "}}"}
	if err := ValidateTemplate(cfg); err == nil {
		t.Error("expected error for whitespace in left_delim")
	}
}

func TestValidateTemplate_EmptyVarKey(t *testing.T) {
	cfg := TemplateConfig{
		Enabled:    true,
		LeftDelim:  "{{",
		RightDelim: "}}",
		Vars:       map[string]string{"": "value"},
	}
	if err := ValidateTemplate(cfg); err == nil {
		t.Error("expected error for empty var key")
	}
}
