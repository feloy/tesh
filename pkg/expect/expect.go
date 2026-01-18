package expect

import (
	"io"

	"github.com/feloy/tesh/pkg/scenarios"
	"mvdan.cc/sh/v3/interp"
)

func CheckExpectations(expectations *scenarios.Expect, result interp.ExitStatus, stdout io.Reader, stderr io.Reader) {
	//fmt.Printf("exit code: %d\n", result)
	//fmt.Printf("stdout: %s\n", stdout)
	//fmt.Printf("stderr: %s\n", stderr)
}
