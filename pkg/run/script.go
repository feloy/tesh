package run

import (
	"context"
	"io"
	"os"

	"github.com/feloy/tesh/pkg/system"
	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

func Script(file io.Reader) {
	script, _ := syntax.NewParser().Parse(file, "")

	runner, _ := interp.New(
		interp.Env(expand.ListEnviron("GLOBAL=global_value")),
		interp.Env(expand.ListEnviron(os.Environ()...)),
		interp.StdIO(nil, os.Stdout, os.Stderr),
	)
	result := runner.Run(context.TODO(), script)

	if status, ok := result.(interp.ExitStatus); ok {
		system.Exit(int(status))
	} else {
		system.Exit(0)
	}
}
