package config

import "testing"

func TestDefaultCacheConfig_Values(t *testing.T) {
	c := DefaultCacheConfig()
	if c.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if c.Dir == "" {
		t.Error("expected Dir to have a default value")
	}
	if c.TTLSeconds <= 0 {
		t.Errorf("expected TTLSeconds > 0, got %d", c.TTLSeconds)
	}
	if c.MaxEntries <= 0 {
		t.Errorf("expected MaxEntries > 0, got %d", c.MaxEntries)
	}
}

func TestValidateCache_DisabledSkipsValidation(t *testing.T) {
	c := CacheConfig{Enabled: false, Dir: "", TTLSeconds: -99, MaxEntries: -1}
	if err := ValidateCache(c); err != nil {
		t.Errorf("expected no error when disabled, got: %v", err)
	}
}

func TestValidateCache_Valid(t *testing.T) {
	c := CacheConfig{Enabled: true, Dir: "/tmp/cache", TTLSeconds: 60, MaxEntries: 50}
	if err := ValidateCache(c); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateCache_EmptyDir(t *testing.T) {
	c := CacheConfig{Enabled: true, Dir: "", TTLSeconds: 60, MaxEntries: 50}
	if err := ValidateCache(c); err == nil {
		t.Error("expected error for empty Dir")
	}
}

func TestValidateCache_NegativeTTL(t *testing.T) {
	c := CacheConfig{Enabled: true, Dir: "/tmp/cache", TTLSeconds: -1, MaxEntries: 50}
	if err := ValidateCache(c); err == nil {
		t.Error("expected error for negative TTLSeconds")
	}
}

func TestValidateCache_NegativeMaxEntries(t *testing.T) {
	c := CacheConfig{Enabled: true, Dir: "/tmp/cache", TTLSeconds: 0, MaxEntries: -5}
	if err := ValidateCache(c); err == nil {
		t.Error("expected error for negative MaxEntries")
	}
}

func TestValidateCache_ZeroTTLAndMaxEntriesMeansUnlimited(t *testing.T) {
	c := CacheConfig{Enabled: true, Dir: "/tmp/cache", TTLSeconds: 0, MaxEntries: 0}
	if err := ValidateCache(c); err != nil {
		t.Errorf("expected no error for zero TTL and MaxEntries, got: %v", err)
	}
}
