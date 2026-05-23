package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Step represents a single runbook step with its shell command and metadata.
type Step struct {
	ID          int
	Title       string
	Command     string
	Rollback    string
	Description string
}

// Runbook holds the parsed representation of a markdown runbook file.
type Runbook struct {
	Title string
	Steps []Step
}

// ParseFile reads a markdown file and extracts structured runbook steps.
func ParseFile(path string) (*Runbook, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open runbook: %w", err)
	}
	defer f.Close()

	rb := &Runbook{}
	var current *Step
	var inCode bool
	var codeTarget *string
	stepID := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "# ") && rb.Title == "" {
			rb.Title = strings.TrimPrefix(line, "# ")
			continue
		}

		if strings.HasPrefix(line, "## ") {
			if current != nil {
				rb.Steps = append(rb.Steps, *current)
			}
			stepID++
			current = &Step{ID: stepID, Title: strings.TrimPrefix(line, "## ")}
			codeTarget = nil
			inCode = false
			continue
		}

		if current == nil {
			continue
		}

		if strings.HasPrefix(line, "<!-- rollback") {
			codeTarget = &current.Rollback
			continue
		}

		if strings.HasPrefix(line, "```") {
			if inCode {
				inCode = false
				codeTarget = nil
			} else {
				inCode = true
				if codeTarget == nil {
					codeTarget = &current.Command
				}
			}
			continue
		}

		if inCode && codeTarget != nil {
			if *codeTarget != "" {
				*codeTarget += "\n"
			}
			*codeTarget += line
			continue
		}

		if !inCode && line != "" {
			if current.Description != "" {
				current.Description += " "
			}
			current.Description += line
		}
	}

	if current != nil {
		rb.Steps = append(rb.Steps, *current)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan runbook: %w", err)
	}

	return rb, nil
}
