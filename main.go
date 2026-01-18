package main

import (
	"context"
	"io"
	"log"
	"os"
	"strings"

	"github.com/feloy/tesh/pkg/handlers/exec"
	"github.com/feloy/tesh/pkg/scenarios"
	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

type nopWriterCloser struct {
	*strings.Reader
}

func (nopWriterCloser) Write([]byte) (int, error) { return 0, io.EOF }
func (nopWriterCloser) Close() error              { return nil }

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	script, _ := syntax.NewParser().Parse(file, "")

	execHandlers := []func(next interp.ExecHandlerFunc) interp.ExecHandlerFunc{}
	if len(os.Args) > 2 {
		mocksFile, err := os.Open(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		defer mocksFile.Close()
		mocksDefinitions, err := scenarios.ParseScenarios(mocksFile)
		if err != nil {
			log.Fatal(err)
		}
		userScenario := os.Args[3]
		for _, scenario := range mocksDefinitions.Scenarios {
			if scenario.ID != userScenario {
				continue
			}
			for _, mock := range scenario.Mocks {
				execHandlers = append(execHandlers, exec.GetExecHandler(mock))
			}
		}
	}

	runner, _ := interp.New(
		interp.Env(expand.ListEnviron("GLOBAL=global_value")),
		interp.Env(expand.ListEnviron(os.Environ()...)),
		interp.StdIO(nil, os.Stdout, os.Stderr),
		interp.ExecHandlers(execHandlers...),
	)
	result := runner.Run(context.TODO(), script)
	if status, ok := result.(interp.ExitStatus); ok {
		os.Exit(int(status))
	} else {
		os.Exit(1)
	}
}
