package run

import (
	"bytes"
	"context"
	"io"
	"log"
	"os"

	"github.com/feloy/tesh/pkg/api"
	"github.com/feloy/tesh/pkg/expect"
	"github.com/feloy/tesh/pkg/handlers/call"
	"github.com/feloy/tesh/pkg/handlers/coverage"
	"github.com/feloy/tesh/pkg/handlers/exec"
	fileHandler "github.com/feloy/tesh/pkg/handlers/file"
	"github.com/feloy/tesh/pkg/output"
	"github.com/feloy/tesh/pkg/scenarios"
	"github.com/feloy/tesh/pkg/system"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

type ScenariosOptions struct {
	SingleScenarioID *string
	WithCoverage     string // empty string for no coverage, "-" for stdout and a filename otherwise
	FilePath         string
}

func Scenarios(file io.Reader, scenariosFile io.Reader, options ScenariosOptions) []api.ScenarioResult {
	script, _ := syntax.NewParser().Parse(file, "")

	mocksDefinitions, err := scenarios.ParseScenarios(scenariosFile)
	if err != nil {
		log.Fatal(err)
	}

	var cov *coverage.Coverage
	if options.WithCoverage != "" {
		cov = coverage.New(script)
	}

	var scenarioResults []api.ScenarioResult
	var found bool = false
	for _, scenario := range mocksDefinitions.Scenarios {
		var expectations *scenarios.Expect

		scenarioEnv := append(os.Environ(), scenario.Envs...)
		var runnerOptions []interp.RunnerOption = []interp.RunnerOption{
			interp.Env(expand.ListEnviron(scenarioEnv...)),
		}

		execHandlers := []func(next interp.ExecHandlerFunc) interp.ExecHandlerFunc{}

		if options.SingleScenarioID != nil && scenario.ID != *options.SingleScenarioID {
			continue
		}
		scenarioResult := api.ScenarioResult{
			ScenarioID: scenario.ID,
		}

		for _, mock := range scenario.Mocks {
			execHandlers = append(execHandlers, exec.GetExecHandler(mock))
		}
		expectations = scenario.Expect
		found = true

		runnerOptions = append(runnerOptions, interp.ExecHandlers(execHandlers...))

		if len(scenario.Files) > 0 {
			runnerOptions = append(runnerOptions, interp.StatHandler(fileHandler.GetStatHandler(scenario.Files)))
		}

		var callsResult call.CallsResult

		var stdout io.ReadWriter = os.Stdout
		var stderr io.ReadWriter = os.Stderr
		if options.WithCoverage == "-" {
			stdout = nil
			stderr = nil
		}
		if expectations != nil {
			var stdoutBuffer bytes.Buffer
			var stderrBuffer bytes.Buffer
			stdout = &stdoutBuffer
			stderr = &stderrBuffer

			if len(expectations.Calls) > 0 {
				var callHandler interp.CallHandlerFunc
				callHandler, callsResult = call.GetCallHandler(expectations.Calls)
				runnerOptions = append(runnerOptions, interp.CallHandler(callHandler))
			}
		}
		runnerOptions = append(runnerOptions, interp.StdIO(nil, stdout, stderr))

		if options.WithCoverage != "" {
			runnerOptions = append(runnerOptions, interp.CallHandler(cov.GetCoverageHandler()))
		}

		runner, _ := interp.New(runnerOptions...)
		result := runner.Run(context.TODO(), script)

		var intResult int
		if _, ok := result.(interp.ExitStatus); ok {
			intResult = int(result.(interp.ExitStatus))
		} else {
			intResult = 0
		}
		if expectations == nil {
			if cov != nil {
				displayCoverage(cov, options)
			}
			system.Exit(intResult)
			return nil
		} else {
			expect.CheckExpectations(expectations, &scenarioResult, intResult, stdout, stderr)
			if len(expectations.Calls) > 0 {
				callsResult.CheckResults(&scenarioResult)
			}
		}
		scenarioResults = append(scenarioResults, scenarioResult)
	}

	if cov != nil {
		displayCoverage(cov, options)
	}

	if !found {
		log.Fatalf("scenario %s not found", *options.SingleScenarioID)
	}
	return scenarioResults
}

func displayCoverage(cov *coverage.Coverage, options ScenariosOptions) {
	if options.WithCoverage != "" {
		positions, lens, covered := cov.GetCoverageResult()
		if options.WithCoverage == "-" {
			scriptFile, err := os.Open(options.FilePath)
			if err != nil {
				log.Fatalf("failed to open file: %v", err)
			}
			defer scriptFile.Close()
			output.OutputCoverage(os.Stdout, scriptFile, positions, lens, covered)
		} else {
			coverageFile, err := os.Create(options.WithCoverage)
			if err != nil {
				log.Fatalf("failed to open file: %v", err)
			}
			defer coverageFile.Close()
			output.OutputCoverageTxt(coverageFile, options.FilePath, positions, lens, covered)
		}
	}
}
