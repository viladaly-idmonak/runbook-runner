package config

import "fmt"

// NotifyOn controls which step outcomes trigger a notification.
type NotifyOn string

const (
	NotifyOnSuccess NotifyOn = "success"
	NotifyOnFailure NotifyOn = "failure"
	NotifyOnAlways  NotifyOn = "always"
)

// NotifyConfig holds configuration for post-run notifications.
type NotifyConfig struct {
	// Enabled toggles the notification system.
	Enabled bool `yaml:"enabled" json:"enabled"`

	// Channel is the notification channel: "webhook", "log".
	Channel string `yaml:"channel" json:"channel"`

	// WebhookURL is the URL to POST results to (used when Channel == "webhook").
	WebhookURL string `yaml:"webhook_url" json:"webhook_url"`

	// On determines when to send a notification.
	On NotifyOn `yaml:"on" json:"on"`
}

// DefaultNotifyConfig returns a safe default notification configuration.
func DefaultNotifyConfig() NotifyConfig {
	return NotifyConfig{
		Enabled: false,
		Channel: "log",
		On:      NotifyOnFailure,
	}
}

var validChannels = map[string]bool{
	"webhook": true,
	"log":     true,
}

var validNotifyOn = map[NotifyOn]bool{
	NotifyOnSuccess: true,
	NotifyOnFailure: true,
	NotifyOnAlways:  true,
}

// ValidateNotify returns an error if the NotifyConfig is invalid.
func ValidateNotify(c NotifyConfig) error {
	if !c.Enabled {
		return nil
	}
	if !validChannels[c.Channel] {
		return fmt.Errorf("notify: unknown channel %q (want \"webhook\" or \"log\")", c.Channel)
	}
	if c.Channel == "webhook" && c.WebhookURL == "" {
		return fmt.Errorf("notify: webhook_url is required when channel is \"webhook\"")
	}
	if !validNotifyOn[c.On] {
		return fmt.Errorf("notify: unknown 'on' value %q (want success|failure|always)", c.On)
	}
	return nil
}
