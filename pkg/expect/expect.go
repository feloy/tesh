package expect

import (
	"io"

	"github.com/feloy/tesh/pkg/api"
	"github.com/feloy/tesh/pkg/scenarios"
)

func CheckExpectations(expectations *scenarios.Expect, scenarioResult *api.ScenarioResult, result int, stdout io.Reader, stderr io.Reader) error {
	if expectations.ExitCode != nil {
		if *expectations.ExitCode != result {
			scenarioResult.ExitCode = &api.ExitCodeError{
				Expected: *expectations.ExitCode,
				Actual:   result,
			}
		}
	}
	if expectations.Stdout != nil {
		stdoutContent, err := io.ReadAll(stdout)
		if err != nil {
			return err
		}
		if *expectations.Stdout != string(stdoutContent) {
			scenarioResult.Stdout = &api.StdioError{
				Expected: *expectations.Stdout,
				Actual:   string(stdoutContent),
			}
		}
	}
	if expectations.Stderr != nil {
		stderrContent, err := io.ReadAll(stderr)
		if err != nil {
			return err
		}
		if *expectations.Stderr != string(stderrContent) {
			scenarioResult.Stderr = &api.StdioError{
				Expected: *expectations.Stderr,
				Actual:   string(stderrContent),
			}
		}
	}
	return nil
}
