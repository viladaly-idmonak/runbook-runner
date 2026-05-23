package executor

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/user/runbook-runner/internal/config"
	"github.com/user/runbook-runner/internal/reporter"
)

var sampleSteps = []reporter.StepEntry{
	{Name: "step-1", Status: reporter.StatusOK},
}

func TestNotifier_DisabledSendsNothing(t *testing.T) {
	var buf bytes.Buffer
	n := NewNotifier(config.NotifyConfig{Enabled: false}, &buf)
	if err := n.Notify("book", sampleSteps, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Error("expected no output for disabled notifier")
	}
}

func TestNotifier_LogChannel_Success(t *testing.T) {
	var buf bytes.Buffer
	cfg := config.NotifyConfig{Enabled: true, Channel: "log", On: config.NotifyOnAlways}
	n := NewNotifier(cfg, &buf)
	if err := n.Notify("my-runbook", sampleSteps, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "SUCCESS") {
		t.Errorf("expected SUCCESS in log output, got %q", buf.String())
	}
}

func TestNotifier_LogChannel_Failure(t *testing.T) {
	var buf bytes.Buffer
	cfg := config.NotifyConfig{Enabled: true, Channel: "log", On: config.NotifyOnAlways}
	n := NewNotifier(cfg, &buf)
	_ = n.Notify("my-runbook", sampleSteps, false)
	if !strings.Contains(buf.String(), "FAILURE") {
		t.Errorf("expected FAILURE in log output, got %q", buf.String())
	}
}

func TestNotifier_OnFailure_SkipsWhenSuccess(t *testing.T) {
	var buf bytes.Buffer
	cfg := config.NotifyConfig{Enabled: true, Channel: "log", On: config.NotifyOnFailure}
	n := NewNotifier(cfg, &buf)
	_ = n.Notify("my-runbook", sampleSteps, true)
	if buf.Len() != 0 {
		t.Error("expected no output when on=failure and run succeeded")
	}
}

func TestNotifier_WebhookChannel_PostsJSON(t *testing.T) {
	var received bytes.Buffer
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = received.ReadFrom(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	cfg := config.NotifyConfig{
		Enabled:    true,
		Channel:    "webhook",
		WebhookURL: ts.URL,
		On:         config.NotifyOnAlways,
	}
	n := NewNotifier(cfg, &bytes.Buffer{})
	if err := n.Notify("hook-book", sampleSteps, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var payload map[string]interface{}
	if err := json.Unmarshal(received.Bytes(), &payload); err != nil {
		t.Fatalf("invalid JSON payload: %v", err)
	}
	if payload["title"] != "hook-book" {
		t.Errorf("expected title \"hook-book\", got %v", payload["title"])
	}
	if payload["status"] != "success" {
		t.Errorf("expected status \"success\", got %v", payload["status"])
	}
}

func TestNotifier_WebhookChannel_Non2xxReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	cfg := config.NotifyConfig{
		Enabled:    true,
		Channel:    "webhook",
		WebhookURL: ts.URL,
		On:         config.NotifyOnAlways,
	}
	n := NewNotifier(cfg, &bytes.Buffer{})
	if err := n.Notify("book", sampleSteps, false); err == nil {
		t.Error("expected error for non-2xx webhook response")
	}
}
