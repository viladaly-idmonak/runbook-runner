package config

import "fmt"

// HookConfig holds the command and lifecycle type for a single hook entry.
type HookConfig struct {
	Type    string `yaml:"type"`
	Command string `yaml:"command"`
}

// validHookTypes lists accepted hook type values.
var validHookTypes = map[string]bool{
	"pre-step":  true,
	"post-step": true,
	"on-error":  true,
}

// ValidateHooks checks that every hook in the slice has a valid type and
// non-empty command.
func ValidateHooks(hooks []HookConfig) error {
	for i, h := range hooks {
		if !validHookTypes[h.Type] {
			return fmt.Errorf("%w: hooks[%d] has unknown type %q", ErrInvalidConfig, i, h.Type)
		}
		if h.Command == "" {
			return fmt.Errorf("%w: hooks[%d] command must not be empty", ErrInvalidConfig, i)
		}
	}
	return nil
}
