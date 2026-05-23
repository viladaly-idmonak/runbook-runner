package executor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/user/runbook-runner/internal/config"
	"github.com/user/runbook-runner/internal/reporter"
)

// Notifier sends a summary notification after a runbook execution.
type Notifier struct {
	cfg    config.NotifyConfig
	client *http.Client
	log    io.Writer
}

// NewNotifier creates a Notifier from the given config, writing log-channel
// output to w.
func NewNotifier(cfg config.NotifyConfig, w io.Writer) *Notifier {
	return &Notifier{
		cfg: cfg,
		client: &http.Client{Timeout: 10 * time.Second},
		log: w,
	}
}

// Notify dispatches a notification based on the run outcome.
// ok is true when all steps succeeded.
func (n *Notifier) Notify(title string, steps []reporter.StepEntry, ok bool) error {
	if !n.cfg.Enabled {
		return nil
	}
	if !n.shouldSend(ok) {
		return nil
	}
	switch n.cfg.Channel {
	case "log":
		return n.sendLog(title, steps, ok)
	case "webhook":
		return n.sendWebhook(title, steps, ok)
	default:
		return fmt.Errorf("notifier: unknown channel %q", n.cfg.Channel)
	}
}

func (n *Notifier) shouldSend(ok bool) bool {
	switch n.cfg.On {
	case config.NotifyOnAlways:
		return true
	case config.NotifyOnSuccess:
		return ok
	case config.NotifyOnFailure:
		return !ok
	}
	return false
}

func (n *Notifier) sendLog(title string, _ []reporter.StepEntry, ok bool) error {
	status := "SUCCESS"
	if !ok {
		status = "FAILURE"
	}
	_, err := fmt.Fprintf(n.log, "[notify] runbook %q finished: %s\n", title, status)
	return err
}

type webhookPayload struct {
	Title  string               `json:"title"`
	Status string               `json:"status"`
	Steps  []reporter.StepEntry `json:"steps"`
}

func (n *Notifier) sendWebhook(title string, steps []reporter.StepEntry, ok bool) error {
	status := "success"
	if !ok {
		status = "failure"
	}
	body, err := json.Marshal(webhookPayload{Title: title, Status: status, Steps: steps})
	if err != nil {
		return fmt.Errorf("notifier: marshal: %w", err)
	}
	resp, err := n.client.Post(n.cfg.WebhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("notifier: post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("notifier: webhook returned status %d", resp.StatusCode)
	}
	return nil
}
