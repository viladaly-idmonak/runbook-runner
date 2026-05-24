package config

import (
	"fmt"
	"os"
)

// SignalConfig controls how the runner responds to OS signals (e.g. SIGINT, SIGTERM).
type SignalConfig struct {
	// GracefulShutdown enables waiting for the current step to finish before exiting.
	GracefulShutdown bool `yaml:"graceful_shutdown"`
	// RunRollbackOnSignal triggers rollback commands when a signal is caught.
	RunRollbackOnSignal bool `yaml:"run_rollback_on_signal"`
	// NotifySignal is the OS signal name to listen for (e.g. "SIGINT", "SIGTERM").
	NotifySignal string `yaml:"notify_signal"`
}

// DefaultSignalConfig returns a SignalConfig with sensible defaults.
func DefaultSignalConfig() SignalConfig {
	return SignalConfig{
		GracefulShutdown:    true,
		RunRollbackOnSignal: false,
		NotifySignal:        "SIGINT",
	}
}

var validSignals = map[string]os.Signal{
	"SIGINT":  os.Interrupt,
	"SIGTERM": nil, // syscall.SIGTERM — kept as nil for portability in this layer
}

// ValidateSignal returns an error if the SignalConfig is invalid.
func ValidateSignal(c SignalConfig) error {
	if c.NotifySignal == "" {
		return fmt.Errorf("%w: notify_signal must not be empty", ErrInvalidConfig)
	}
	if _, ok := validSignals[c.NotifySignal]; !ok {
		return fmt.Errorf("%w: unknown notify_signal %q (supported: SIGINT, SIGTERM)", ErrInvalidConfig, c.NotifySignal)
	}
	return nil
}
