package config

import "fmt"

// MatrixConfig controls step matrix expansion — running a step multiple times
// with different variable substitutions.
type MatrixConfig struct {
	// Enabled turns matrix expansion on or off.
	Enabled bool `yaml:"enabled"`

	// Vars maps a variable name to a list of values. Each combination of
	// values across all variables produces one expanded step instance.
	// Example: {"ENV": ["staging", "prod"], "REGION": ["us", "eu"]}
	Vars map[string][]string `yaml:"vars"`

	// MaxExpansions is the upper bound on the total number of expanded
	// instances. 0 means no limit.
	MaxExpansions int `yaml:"max_expansions"`
}

// DefaultMatrixConfig returns a MatrixConfig with sensible defaults.
func DefaultMatrixConfig() MatrixConfig {
	return MatrixConfig{
		Enabled:       false,
		Vars:          map[string][]string{},
		MaxExpansions: 50,
	}
}

// ValidateMatrix returns an error if cfg contains invalid settings.
func ValidateMatrix(cfg MatrixConfig) error {
	if !cfg.Enabled {
		return nil
	}

	if len(cfg.Vars) == 0 {
		return fmt.Errorf("matrix: enabled but no vars defined")
	}

	for key, values := range cfg.Vars {
		if key == "" {
			return fmt.Errorf("matrix: var key must not be empty")
		}
		if len(values) == 0 {
			return fmt.Errorf("matrix: var %q has no values", key)
		}
		for _, v := range values {
			if v == "" {
				return fmt.Errorf("matrix: var %q contains an empty value", key)
			}
		}
	}

	if cfg.MaxExpansions < 0 {
		return fmt.Errorf("matrix: max_expansions must be >= 0, got %d", cfg.MaxExpansions)
	}

	total := expansionCount(cfg.Vars)
	if cfg.MaxExpansions > 0 && total > cfg.MaxExpansions {
		return fmt.Errorf("matrix: expansion count %d exceeds max_expansions %d", total, cfg.MaxExpansions)
	}

	return nil
}

// expansionCount returns the Cartesian product size of all var value lists.
func expansionCount(vars map[string][]string) int {
	count := 1
	for _, values := range vars {
		count *= len(values)
	}
	return count
}
