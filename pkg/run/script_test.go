package run

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/feloy/tesh/pkg/system"
)

func TestScript(t *testing.T) {
	exitCalled := false
	exitCode := 0

	system.Exit = func(code int) {
		exitCalled = true
		exitCode = code
	}

	script := strings.NewReader(`echo -n /$MYVAR/; >&2 echo "error"`)
	os.Setenv("MYVAR", "myvalue")

	oldStdout := os.Stdout
	oldStderr := os.Stderr
	rStdout, wStdout, _ := os.Pipe()
	rStderr, wStderr, _ := os.Pipe()
	os.Stdout = wStdout
	os.Stderr = wStderr
	outCStdout := make(chan string)
	outCStderr := make(chan string)

	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, rStdout)
		outCStdout <- buf.String()
	}()

	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, rStderr)
		outCStderr <- buf.String()
	}()

	Script(script, ScriptOptions{WithCoverage: "", FilePath: ""})
	wStdout.Close()
	wStderr.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr
	outStdout := <-outCStdout
	outStderr := <-outCStderr

	if !exitCalled {
		t.Fatalf("os.Exit was not called")
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	if outStdout != "/myvalue/" {
		t.Fatalf("expected stdout to be /myvalue/, got %s", outStdout)
	}
	if outStderr != "error\n" {
		t.Fatalf("expected stderr to be error, got %s", outStderr)
	}
}

func TestScriptWithExitCode(t *testing.T) {
	exitCalled := false
	exitCode := 0

	system.Exit = func(code int) {
		exitCalled = true
		exitCode = code
	}

	script := strings.NewReader(`exit 126`)
	Script(script, ScriptOptions{WithCoverage: "", FilePath: ""})
	if !exitCalled {
		t.Fatalf("os.Exit was not called")
	}
	if exitCode != 126 {
		t.Fatalf("expected exit code 126, got %d", exitCode)
	}
}
