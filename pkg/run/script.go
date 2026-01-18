package run

import (
	"context"
	"os"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

func Script(file *os.File) {
	script, _ := syntax.NewParser().Parse(file, "")

	runner, _ := interp.New(
		interp.Env(expand.ListEnviron("GLOBAL=global_value")),
		interp.Env(expand.ListEnviron(os.Environ()...)),
		interp.StdIO(nil, os.Stdout, os.Stderr),
	)
	result := runner.Run(context.TODO(), script)

	if status, ok := result.(interp.ExitStatus); ok {
		os.Exit(int(status))
	} else {
		os.Exit(1)
	}
}
