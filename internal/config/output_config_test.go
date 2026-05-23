package config

import "testing"

func TestDefaultOutputConfig_Values(t *testing.T) {
	c := DefaultOutputConfig()
	if c.Format != OutputFormatText {
		t.Errorf("expected format %q, got %q", OutputFormatText, c.Format)
	}
	if c.Verbose {
		t.Error("expected Verbose to be false")
	}
	if !c.Color {
		t.Error("expected Color to be true")
	}
	if c.TimestampsEnabled {
		t.Error("expected TimestampsEnabled to be false")
	}
}

func TestValidateOutput_ValidText(t *testing.T) {
	if err := ValidateOutput(OutputConfig{Format: OutputFormatText}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateOutput_ValidJSON(t *testing.T) {
	if err := ValidateOutput(OutputConfig{Format: OutputFormatJSON}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateOutput_EmptyFormat(t *testing.T) {
	err := ValidateOutput(OutputConfig{Format: ""})
	if err == nil {
		t.Fatal("expected error for empty format")
	}
}

func TestValidateOutput_UnknownFormat(t *testing.T) {
	err := ValidateOutput(OutputConfig{Format: "xml"})
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
	if want := "unknown output format"; !containsStr(err.Error(), want) {
		t.Errorf("error %q should contain %q", err.Error(), want)
	}
}

func TestValidateOutput_DefaultIsValid(t *testing.T) {
	if err := ValidateOutput(DefaultOutputConfig()); err != nil {
		t.Errorf("default config should be valid, got: %v", err)
	}
}

// containsStr is a small helper shared across config tests.
func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		func() bool {
			for i := 0; i+len(sub) <= len(s); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}())
}
