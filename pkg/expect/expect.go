package expect

import (
	"io"
	"log"

	"github.com/feloy/tesh/pkg/scenarios"
)

func CheckExpectations(expectations *scenarios.Expect, result int, stdout io.Reader, stderr io.Reader) {
	if expectations.ExitCode != nil {
		if *expectations.ExitCode != result {
			log.Fatalf("expected exit code %d, got %d", *expectations.ExitCode, result)
		}
		if expectations.Stdout != nil {
			stdoutContent, err := io.ReadAll(stdout)
			if err != nil {
				log.Fatalf("failed to read stdout: %v", err)
			}
			if *expectations.Stdout != string(stdoutContent) {
				log.Fatalf("expected stdout %q, got %q", *expectations.Stdout, string(stdoutContent))
			}
		}
		if expectations.Stderr != nil {
			stderrContent, err := io.ReadAll(stderr)
			if err != nil {
				log.Fatalf("failed to read stderr: %v", err)
			}
			if *expectations.Stderr != string(stderrContent) {
				log.Fatalf("expected stderr %q, got %q", *expectations.Stderr, string(stderrContent))
			}
		}
	}
}
