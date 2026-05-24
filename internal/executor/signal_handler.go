package executor

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/your-org/runbook-runner/internal/config"
)

// SignalHandler listens for OS signals and cancels the run context.
type SignalHandler struct {
	cfg    config.SignalConfig
	cancel context.CancelFunc
	ch     chan os.Signal
}

// NewSignalHandler creates a SignalHandler and begins listening immediately.
// Call Stop() to clean up the goroutine.
func NewSignalHandler(cfg config.SignalConfig, cancel context.CancelFunc) *SignalHandler {
	sh := &SignalHandler{
		cfg:    cfg,
		cancel: cancel,
		ch:     make(chan os.Signal, 1),
	}
	sig := resolveSignal(cfg.NotifySignal)
	signal.Notify(sh.ch, sig)
	go sh.listen()
	return sh
}

// Stop unregisters the signal listener and closes the internal channel.
func (sh *SignalHandler) Stop() {
	signal.Stop(sh.ch)
	close(sh.ch)
}

func (sh *SignalHandler) listen() {
	for s := range sh.ch {
		log.Printf("[signal] received %s — requesting shutdown", s)
		sh.cancel()
		return
	}
}

func resolveSignal(name string) os.Signal {
	switch name {
	case "SIGTERM":
		return syscall.SIGTERM
	default:
		return os.Interrupt // SIGINT
	}
}
