package mocks

import (
	"strings"
	"testing"
)

func TestParseMocks(t *testing.T) {
	mocks, err := ParseMocks(strings.NewReader(`
scenarios:
- id: no-change-in-repository
  description: there is no change in the repository
  mocks:
  - description: git diff indicates there is no change in the repository
    command: git
    args:
    - diff
    exit-code: 0
    stdout: some text in standard output
    stderr: some text in standard error
`))
	if err != nil {
		t.Fatalf("failed to parse mocks: %v", err)
	}
	if len(mocks.Scenarios) != 1 {
		t.Fatalf("expected 1 scenario, got %d", len(mocks.Scenarios))
	}
	if mocks.Scenarios[0].ID != "no-change-in-repository" {
		t.Fatalf("expected scenario ID to be no-change-in-repository, got %s", mocks.Scenarios[0].ID)
	}
	if mocks.Scenarios[0].Mocks[0].Description != "git diff indicates there is no change in the repository" {
		t.Fatalf("expected mock description to be git diff indicates there is no change in the repository, got %s", mocks.Scenarios[0].Mocks[0].Description)
	}
	if mocks.Scenarios[0].Mocks[0].Command != "git" {
		t.Fatalf("expected mock command to be git, got %s", mocks.Scenarios[0].Mocks[0].Command)
	}
	if len(mocks.Scenarios[0].Mocks[0].Args) != 1 {
		t.Fatalf("expected mock args to be diff, got %d", len(mocks.Scenarios[0].Mocks[0].Args))
	}
	if mocks.Scenarios[0].Mocks[0].Args[0] != "diff" {
		t.Fatalf("expected mock args to be diff, got %s", mocks.Scenarios[0].Mocks[0].Args[0])
	}
	if mocks.Scenarios[0].Mocks[0].ExitCode == nil {
		t.Fatalf("expected mock exit code to be 0, got nil")
	} else if *mocks.Scenarios[0].Mocks[0].ExitCode != 0 {
		t.Fatalf("expected mock exit code to be 0, got %d", *mocks.Scenarios[0].Mocks[0].ExitCode)
	}
	if mocks.Scenarios[0].Mocks[0].Stdout == nil {
		t.Fatalf("expected mock stdout to be some text in standard output, got nil")
	} else if *mocks.Scenarios[0].Mocks[0].Stdout != "some text in standard output" {
		t.Fatalf("expected mock stdout to be some text in standard output, got %s", *mocks.Scenarios[0].Mocks[0].Stdout)
	}

	if mocks.Scenarios[0].Mocks[0].Stderr == nil {
		t.Fatalf("expected mock stderr to be some text in standard error, got nil")
	} else if *mocks.Scenarios[0].Mocks[0].Stderr != "some text in standard error" {
		t.Fatalf("expected mock stderr to be some text in standard error, got %s", *mocks.Scenarios[0].Mocks[0].Stderr)
	}
}
