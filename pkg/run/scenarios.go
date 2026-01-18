package run

import (
	"bytes"
	"context"
	"io"
	"log"
	"os"

	"github.com/feloy/tesh/pkg/expect"
	"github.com/feloy/tesh/pkg/handlers/exec"
	"github.com/feloy/tesh/pkg/scenarios"
	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

func Scenarios(file *os.File, scenariosFile *os.File, singleScenarioID *string) {
	if singleScenarioID == nil {
		log.Fatal("single scenario ID is required")
	}

	script, _ := syntax.NewParser().Parse(file, "")

	var expectations *scenarios.Expect

	execHandlers := []func(next interp.ExecHandlerFunc) interp.ExecHandlerFunc{}

	defer scenariosFile.Close()
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

	var stdout io.ReadWriter = os.Stdout
	var stderr io.ReadWriter = os.Stderr
	if expectations != nil {
		var stdoutBuffer bytes.Buffer
		var stderrBuffer bytes.Buffer
		stdout = &stdoutBuffer
		stderr = &stderrBuffer
	}

	runner, _ := interp.New(
		interp.Env(expand.ListEnviron("GLOBAL=global_value")),
		interp.Env(expand.ListEnviron(os.Environ()...)),
		interp.StdIO(nil, stdout, stderr),
		interp.ExecHandlers(execHandlers...),
	)
	result := runner.Run(context.TODO(), script)

	if expectations == nil {
		if status, ok := result.(interp.ExitStatus); ok {
			os.Exit(int(status))
		} else {
			os.Exit(1)
		}
	} else {
		expect.CheckExpectations(expectations, result.(interp.ExitStatus), stdout, stderr)
	}
}
