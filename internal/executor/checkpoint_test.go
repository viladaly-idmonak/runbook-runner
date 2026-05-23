package executor_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/runbook-runner/internal/executor"
)

func TestCheckpointStore_MarkAndIsDone(t *testing.T) {
	dir := t.TempDir()
	store := executor.NewCheckpointStore(dir, "my-runbook")

	if store.IsDone("step-1") {
		t.Fatal("step-1 should not be done yet")
	}
	if err := store.MarkDone("step-1"); err != nil {
		t.Fatalf("MarkDone: %v", err)
	}
	if !store.IsDone("step-1") {
		t.Error("step-1 should be done after MarkDone")
	}
}

func TestCheckpointStore_MarkDone_Idempotent(t *testing.T) {
	dir := t.TempDir()
	store := executor.NewCheckpointStore(dir, "rb")

	for i := 0; i < 3; i++ {
		if err := store.MarkDone("step-a"); err != nil {
			t.Fatalf("MarkDone iteration %d: %v", i, err)
		}
	}
	// File should exist and contain only one entry.
	files, _ := filepath.Glob(filepath.Join(dir, "*.json"))
	if len(files) != 1 {
		t.Fatalf("expected 1 checkpoint file, got %d", len(files))
	}
}

func TestCheckpointStore_Reset(t *testing.T) {
	dir := t.TempDir()
	store := executor.NewCheckpointStore(dir, "rb")

	_ = store.MarkDone("step-1")
	if err := store.Reset(); err != nil {
		t.Fatalf("Reset: %v", err)
	}
	if store.IsDone("step-1") {
		t.Error("step-1 should not be done after Reset")
	}
}

func TestCheckpointStore_Reset_NoFile(t *testing.T) {
	dir := t.TempDir()
	store := executor.NewCheckpointStore(dir, "no-file")
	if err := store.Reset(); err != nil {
		t.Errorf("Reset on missing file should not error, got: %v", err)
	}
}

func TestCheckpointStore_MkdirOnSave(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "dir")
	store := executor.NewCheckpointStore(dir, "rb")
	if err := store.MarkDone("step-1"); err != nil {
		t.Fatalf("MarkDone with nested dir: %v", err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Errorf("expected dir to be created: %v", err)
	}
}
