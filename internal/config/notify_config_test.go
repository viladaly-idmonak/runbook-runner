package config

import "testing"

func TestDefaultNotifyConfig_Values(t *testing.T) {
	c := DefaultNotifyConfig()
	if c.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if c.Channel != "log" {
		t.Errorf("expected Channel \"log\", got %q", c.Channel)
	}
	if c.On != NotifyOnFailure {
		t.Errorf("expected On \"failure\", got %q", c.On)
	}
}

func TestValidateNotify_DisabledSkipsValidation(t *testing.T) {
	c := NotifyConfig{Enabled: false, Channel: "", WebhookURL: "", On: ""}
	if err := ValidateNotify(c); err != nil {
		t.Errorf("expected no error for disabled config, got %v", err)
	}
}

func TestValidateNotify_ValidLogChannel(t *testing.T) {
	c := NotifyConfig{Enabled: true, Channel: "log", On: NotifyOnAlways}
	if err := ValidateNotify(c); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateNotify_ValidWebhookChannel(t *testing.T) {
	c := NotifyConfig{Enabled: true, Channel: "webhook", WebhookURL: "https://example.com/hook", On: NotifyOnSuccess}
	if err := ValidateNotify(c); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateNotify_WebhookMissingURL(t *testing.T) {
	c := NotifyConfig{Enabled: true, Channel: "webhook", On: NotifyOnFailure}
	if err := ValidateNotify(c); err == nil {
		t.Error("expected error for webhook without URL")
	}
}

func TestValidateNotify_UnknownChannel(t *testing.T) {
	c := NotifyConfig{Enabled: true, Channel: "slack", On: NotifyOnAlways}
	if err := ValidateNotify(c); err == nil {
		t.Error("expected error for unknown channel")
	}
}

func TestValidateNotify_UnknownOnValue(t *testing.T) {
	c := NotifyConfig{Enabled: true, Channel: "log", On: "never"}
	if err := ValidateNotify(c); err == nil {
		t.Error("expected error for unknown 'on' value")
	}
}
