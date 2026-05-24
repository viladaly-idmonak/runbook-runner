package config

import (
	"errors"
	"strings"
)

// AuditConfig controls whether step execution events are written to an audit log.
type AuditConfig struct {
	// Enabled turns audit logging on or off.
	Enabled bool `yaml:"enabled"`

	// FilePath is the destination file for audit entries. Required when Enabled.
	FilePath string `yaml:"file_path"`

	// Format is either "text" or "json".
	Format string `yaml:"format"`

	// IncludeOutput appends the command stdout/stderr to each audit entry.
	IncludeOutput bool `yaml:"include_output"`
}

// DefaultAuditConfig returns a safe default configuration.
func DefaultAuditConfig() AuditConfig {
	return AuditConfig{
		Enabled:       false,
		FilePath:      "",
		Format:        "text",
		IncludeOutput: false,
	}
}

// ValidateAudit returns an error if the AuditConfig is invalid.
func ValidateAudit(cfg AuditConfig) error {
	if !cfg.Enabled {
		return nil
	}

	if strings.TrimSpace(cfg.FilePath) == "" {
		return errors.New("audit: file_path must not be empty when audit logging is enabled")
	}

	switch cfg.Format {
	case "text", "json":
		// valid
	case "":
		return errors.New("audit: format must not be empty")
	default:
		return errors.New("audit: format must be \"text\" or \"json\", got \"" + cfg.Format + "\"")
	}

	return nil
}

// IsJSON reports whether the audit log format is set to JSON.
// This is a convenience helper for code that needs to select an encoder.
func (a AuditConfig) IsJSON() bool {
	return strings.EqualFold(a.Format, "json")
}
