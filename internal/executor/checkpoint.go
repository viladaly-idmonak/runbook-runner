package executor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// checkpointState is persisted to disk after each successful step.
type checkpointState struct {
	CompletedSteps []string `json:"completed_steps"`
}

// CheckpointStore reads and writes checkpoint state for a named runbook.
type CheckpointStore struct {
	dir      string
	runbook  string
}

// NewCheckpointStore creates a CheckpointStore that uses dir for storage.
func NewCheckpointStore(dir, runbook string) *CheckpointStore {
	return &CheckpointStore{dir: dir, runbook: runbook}
}

func (s *CheckpointStore) filePath() string {
	safe := strings.NewReplacer("/", "_", " ", "_").Replace(s.runbook)
	return filepath.Join(s.dir, safe+".checkpoint.json")
}

// MarkDone records stepID as successfully completed.
func (s *CheckpointStore) MarkDone(stepID string) error {
	state, _ := s.load()
	for _, id := range state.CompletedSteps {
		if id == stepID {
			return nil
		}
	}
	state.CompletedSteps = append(state.CompletedSteps, stepID)
	return s.save(state)
}

// IsDone reports whether stepID has already been completed.
func (s *CheckpointStore) IsDone(stepID string) bool {
	state, err := s.load()
	if err != nil {
		return false
	}
	for _, id := range state.CompletedSteps {
		if id == stepID {
			return true
		}
	}
	return false
}

// Reset removes the checkpoint file so the runbook starts fresh.
func (s *CheckpointStore) Reset() error {
	err := os.Remove(s.filePath())
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func (s *CheckpointStore) load() (checkpointState, error) {
	var state checkpointState
	data, err := os.ReadFile(s.filePath())
	if err != nil {
		return state, err
	}
	err = json.Unmarshal(data, &state)
	return state, err
}

func (s *CheckpointStore) save(state checkpointState) error {
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return fmt.Errorf("checkpoint: mkdir %s: %w", s.dir, err)
	}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath(), data, 0o644)
}
