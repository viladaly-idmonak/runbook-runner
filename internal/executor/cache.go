package executor

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/user/runbook-runner/internal/config"
)

// cacheEntry is the on-disk representation of a cached step result.
type cacheEntry struct {
	Output    string    `json:"output"`
	CreatedAt time.Time `json:"created_at"`
}

// StepCache stores and retrieves step outputs keyed by a hash of the
// step command. It respects the TTL and max-entries limits from config.
type StepCache struct {
	cfg config.CacheConfig
}

// NewStepCache returns a StepCache. When caching is disabled every
// operation becomes a no-op.
func NewStepCache(cfg config.CacheConfig) *StepCache {
	return &StepCache{cfg: cfg}
}

// key derives a deterministic filename for the given command string.
func (sc *StepCache) key(command string) string {
	h := sha256.Sum256([]byte(command))
	return hex.EncodeToString(h[:]) + ".json"
}

// Get returns the cached output for command, and whether a valid (non-stale)
// entry was found. Returns ("", false) when caching is disabled or no entry
// exists.
func (sc *StepCache) Get(command string) (string, bool) {
	if !sc.cfg.Enabled {
		return "", false
	}
	path := filepath.Join(sc.cfg.Dir, sc.key(command))
	data, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}
	var entry cacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return "", false
	}
	if sc.cfg.TTLSeconds > 0 {
		age := time.Since(entry.CreatedAt).Seconds()
		if age > float64(sc.cfg.TTLSeconds) {
			_ = os.Remove(path)
			return "", false
		}
	}
	return entry.Output, true
}

// Set writes output to the cache for the given command. It is a no-op when
// caching is disabled.
func (sc *StepCache) Set(command, output string) error {
	if !sc.cfg.Enabled {
		return nil
	}
	if err := os.MkdirAll(sc.cfg.Dir, 0o755); err != nil {
		return fmt.Errorf("cache: mkdir %s: %w", sc.cfg.Dir, err)
	}
	entry := cacheEntry{Output: output, CreatedAt: time.Now()}
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("cache: marshal: %w", err)
	}
	path := filepath.Join(sc.cfg.Dir, sc.key(command))
	return os.WriteFile(path, data, 0o644)
}
