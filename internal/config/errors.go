package config

import "errors"

// Sentinel errors returned by Config.Validate.
var (
	ErrEmptyShellPath     = errors.New("config: shell_path must not be empty")
	ErrInvalidOutputFormat = errors.New("config: output_format must be \"text\" or \"json\"")
	ErrInvalidTimeout     = errors.New("config: step_timeout must be greater than zero")
)
