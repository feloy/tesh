package scenarios

import (
	"strings"
	"testing"
)

func TestParseMocks(t *testing.T) {
	scenarios, err := ParseScenarios(strings.NewReader(`
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
	if len(scenarios.Scenarios) != 1 {
		t.Fatalf("expected 1 scenario, got %d", len(scenarios.Scenarios))
	}
	if scenarios.Scenarios[0].ID != "no-change-in-repository" {
		t.Fatalf("expected scenario ID to be no-change-in-repository, got %s", scenarios.Scenarios[0].ID)
	}
	if scenarios.Scenarios[0].Mocks[0].Description != "git diff indicates there is no change in the repository" {
		t.Fatalf("expected mock description to be git diff indicates there is no change in the repository, got %s", scenarios.Scenarios[0].Mocks[0].Description)
	}
	if scenarios.Scenarios[0].Mocks[0].Command != "git" {
		t.Fatalf("expected mock command to be git, got %s", scenarios.Scenarios[0].Mocks[0].Command)
	}
	if len(scenarios.Scenarios[0].Mocks[0].Args) != 1 {
		t.Fatalf("expected mock args to be diff, got %d", len(scenarios.Scenarios[0].Mocks[0].Args))
	}
	if scenarios.Scenarios[0].Mocks[0].Args[0] != "diff" {
		t.Fatalf("expected mock args to be diff, got %s", scenarios.Scenarios[0].Mocks[0].Args[0])
	}
	if scenarios.Scenarios[0].Mocks[0].ExitCode == nil {
		t.Fatalf("expected mock exit code to be 0, got nil")
	} else if *scenarios.Scenarios[0].Mocks[0].ExitCode != 0 {
		t.Fatalf("expected mock exit code to be 0, got %d", *scenarios.Scenarios[0].Mocks[0].ExitCode)
	}
	if scenarios.Scenarios[0].Mocks[0].Stdout == nil {
		t.Fatalf("expected mock stdout to be some text in standard output, got nil")
	} else if *scenarios.Scenarios[0].Mocks[0].Stdout != "some text in standard output" {
		t.Fatalf("expected mock stdout to be some text in standard output, got %s", *scenarios.Scenarios[0].Mocks[0].Stdout)
	}

	if scenarios.Scenarios[0].Mocks[0].Stderr == nil {
		t.Fatalf("expected mock stderr to be some text in standard error, got nil")
	} else if *scenarios.Scenarios[0].Mocks[0].Stderr != "some text in standard error" {
		t.Fatalf("expected mock stderr to be some text in standard error, got %s", *scenarios.Scenarios[0].Mocks[0].Stderr)
	}
}
