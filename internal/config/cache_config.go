package config

import "fmt"

// CacheConfig controls step-output caching behaviour.
type CacheConfig struct {
	// Enabled turns caching on or off.
	Enabled bool `yaml:"enabled"`

	// Dir is the directory where cache entries are stored.
	Dir string `yaml:"dir"`

	// TTLSeconds is the maximum age of a cache entry before it is
	// considered stale. 0 means entries never expire.
	TTLSeconds int `yaml:"ttl_seconds"`

	// MaxEntries caps the total number of cached step outputs kept on
	// disk. 0 means unlimited.
	MaxEntries int `yaml:"max_entries"`
}

// DefaultCacheConfig returns a safe, disabled-by-default configuration.
func DefaultCacheConfig() CacheConfig {
	return CacheConfig{
		Enabled:    false,
		Dir:        ".runbook-cache",
		TTLSeconds: 3600,
		MaxEntries: 100,
	}
}

// ValidateCache returns an error when the configuration is inconsistent.
func ValidateCache(c CacheConfig) error {
	if !c.Enabled {
		return nil
	}
	if c.Dir == "" {
		return fmt.Errorf("cache: dir must not be empty when caching is enabled")
	}
	if c.TTLSeconds < 0 {
		return fmt.Errorf("cache: ttl_seconds must be >= 0, got %d", c.TTLSeconds)
	}
	if c.MaxEntries < 0 {
		return fmt.Errorf("cache: max_entries must be >= 0, got %d", c.MaxEntries)
	}
	return nil
}
