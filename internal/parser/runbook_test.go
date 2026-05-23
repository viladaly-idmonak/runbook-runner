package parser

import (
	"os"
	"testing"
)

const sampleRunbook = `# Deploy Service

## Stop old containers

Stop any running containers before deploying.

` + "```" + `sh
docker stop myapp
` + "```" + `

<!-- rollback -->
` + "```" + `sh
docker start myapp
` + "```" + `

## Pull latest image

Fetch the newest image from the registry.

` + "```" + `sh
docker pull registry.example.com/myapp:latest
` + "```" + `
`

func writeTempRunbook(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "runbook-*.md")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestParseFile_Title(t *testing.T) {
	path := writeTempRunbook(t, sampleRunbook)
	rb, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rb.Title != "Deploy Service" {
		t.Errorf("expected title 'Deploy Service', got %q", rb.Title)
	}
}

func TestParseFile_StepCount(t *testing.T) {
	path := writeTempRunbook(t, sampleRunbook)
	rb, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rb.Steps) != 2 {
		t.Errorf("expected 2 steps, got %d", len(rb.Steps))
	}
}

func TestParseFile_CommandAndRollback(t *testing.T) {
	path := writeTempRunbook(t, sampleRunbook)
	rb, err := ParseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	step := rb.Steps[0]
	if step.Command != "docker stop myapp" {
		t.Errorf("unexpected command: %q", step.Command)
	}
	if step.Rollback != "docker start myapp" {
		t.Errorf("unexpected rollback: %q", step.Rollback)
	}
}

func TestParseFile_MissingFile(t *testing.T) {
	_, err := ParseFile("/nonexistent/runbook.md")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
