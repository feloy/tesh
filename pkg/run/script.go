package run

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/feloy/tesh/pkg/handlers/coverage"
	"github.com/feloy/tesh/pkg/output"
	"github.com/feloy/tesh/pkg/system"
	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

type ScriptOptions struct {
	WithCoverage string // empty string for no coverage, "-" for stdout and a filename otherwise
	FilePath     string
}

func Script(file io.Reader, options ScriptOptions) {
	script, _ := syntax.NewParser().Parse(file, "")

	var stdout io.ReadWriter = os.Stdout
	var stderr io.ReadWriter = os.Stderr
	// coverage suppresses stdout and stderr and displays covered lines in stdout
	if options.WithCoverage == "-" {
		stdout = nil
		stderr = nil
	}

	runnerOptions := []interp.RunnerOption{
		interp.Env(expand.ListEnviron(os.Environ()...)),
		interp.StdIO(nil, stdout, stderr),
	}

	var cov *coverage.Coverage
	if options.WithCoverage != "" {
		cov = coverage.New(script)
		runnerOptions = append(runnerOptions, interp.CallHandler(cov.GetCoverageHandler()))
	}

	runner, _ := interp.New(runnerOptions...)
	result := runner.Run(context.TODO(), script)

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
	if status, ok := result.(interp.ExitStatus); ok {
		system.Exit(int(status))
	} else {
		system.Exit(0)
	}
}
