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
	"github.com/feloy/tesh/pkg/handlers/exec"
	"github.com/feloy/tesh/pkg/scenarios"
	"github.com/feloy/tesh/pkg/system"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

func Scenarios(file io.Reader, scenariosFile io.Reader, singleScenarioID *string) []api.ScenarioResult {
	if singleScenarioID == nil {
		log.Fatal("single scenario ID is required")
	}

	var scenarioResult = api.ScenarioResult{
		ScenarioID: *singleScenarioID,
	}

	script, _ := syntax.NewParser().Parse(file, "")

	var expectations *scenarios.Expect

	var runnerOptions []interp.RunnerOption = []interp.RunnerOption{
		interp.Env(expand.ListEnviron("GLOBAL=global_value")),
		interp.Env(expand.ListEnviron(os.Environ()...)),
	}

	execHandlers := []func(next interp.ExecHandlerFunc) interp.ExecHandlerFunc{}

	mocksDefinitions, err := scenarios.ParseScenarios(scenariosFile)
	if err != nil {
		log.Fatal(err)
	}
	var found bool = false
	for _, scenario := range mocksDefinitions.Scenarios {
		if scenario.ID != *singleScenarioID {
			continue
		}
		for _, mock := range scenario.Mocks {
			execHandlers = append(execHandlers, exec.GetExecHandler(mock))
		}
		expectations = scenario.Expect
		found = true
		break
	}

	if !found {
		log.Fatalf("scenario %s not found", *singleScenarioID)
	}

	runnerOptions = append(runnerOptions, interp.ExecHandlers(execHandlers...))

	var callsResult call.CallsResult

	var stdout io.ReadWriter = os.Stdout
	var stderr io.ReadWriter = os.Stderr
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

	runner, _ := interp.New(runnerOptions...)
	result := runner.Run(context.TODO(), script)

	var intResult int
	if _, ok := result.(interp.ExitStatus); ok {
		intResult = int(result.(interp.ExitStatus))
	} else {
		intResult = 0
	}
	if expectations == nil {
		system.Exit(intResult)
		return nil
	} else {
		expect.CheckExpectations(expectations, &scenarioResult, intResult, stdout, stderr)
		if len(expectations.Calls) > 0 {
			callsResult.CheckResults(&scenarioResult)
		}
		return []api.ScenarioResult{scenarioResult}
	}
}
