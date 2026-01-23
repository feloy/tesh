package run

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/feloy/tesh/pkg/system"
)

func TestScenarioStdIO(t *testing.T) {
	exitCalled := false
	exitCode := 0

	system.Exit = func(code int) {
		exitCalled = true
		exitCode = code
	}

	script := strings.NewReader(`cat /path/to/file; >&2 echo "error"`)

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
	scenarios := strings.NewReader(`
scenarios:
- id: file-exists
  description: file exists
  mocks:
  - description: the file /path/to/file exists
    command: cat
    args:
    - /path/to/file
    exit-code: 0
    stdout: some text in the file`)
	singleScenarioID := "file-exists"

	Scenarios(script, scenarios, &singleScenarioID)

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
	if outStdout != "some text in the file" {
		t.Fatalf("expected stdout to be 'some text in the file', got %s", outStdout)
	}
	if outStderr != "error\n" {
		t.Fatalf("expected stderr to be error, got %s", outStderr)
	}
}

func TestScenarioEnvs(t *testing.T) {
	system.Exit = func(code int) {}

	script := strings.NewReader(`echo -n /$MYVAR/`)
	os.Setenv("MYVAR", "myvalue")

	oldStdout := os.Stdout
	rStdout, wStdout, _ := os.Pipe()
	os.Stdout = wStdout
	outCStdout := make(chan string)

	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, rStdout)
		outCStdout <- buf.String()
	}()

	Script(script)
	wStdout.Close()
	os.Stdout = oldStdout
	outStdout := <-outCStdout

	if outStdout != "/myvalue/" {
		t.Fatalf("expected stdout to be /myvalue/, got %s", outStdout)
	}
}

func TestScenariosWithExitCode0(t *testing.T) {
	exitCalled := false
	exitCode := 0

	system.Exit = func(code int) {
		exitCalled = true
		exitCode = code
	}

	script := strings.NewReader(`cat /path/to/file`)
	scenarios := strings.NewReader(`
scenarios:
- id: file-exists
  description: file exists
  mocks:
  - description: the file /path/to/file exists
    command: cat
    args:
    - /path/to/file
    exit-code: 0
    stdout: some text in the file`)

	singleScenarioID := "file-exists"
	Scenarios(script, scenarios, &singleScenarioID)

	if !exitCalled {
		t.Fatalf("os.Exit was not called")
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
}

func TestScenariosWithExitCode1(t *testing.T) {
	exitCalled := false
	exitCode := 0

	system.Exit = func(code int) {
		exitCalled = true
		exitCode = code
	}

	script := strings.NewReader(`cat /path/to/file`)
	scenarios := strings.NewReader(`
scenarios:
- id: file-not-exists
  description: file does not exist
  mocks:
  - description: the file /path/to/file exists
    command: cat
    args:
    - /path/to/file
    exit-code: 1
    stdout: some text in the file`)

	singleScenarioID := "file-not-exists"
	Scenarios(script, scenarios, &singleScenarioID)

	if !exitCalled {
		t.Fatalf("os.Exit was not called")
	}
	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}
}

func TestScenariosWithDefaultExit(t *testing.T) {
	exitCalled := false
	exitCode := 0

	system.Exit = func(code int) {
		exitCalled = true
		exitCode = code
	}

	script := strings.NewReader(`echo 1`)
	scenarios := strings.NewReader(`
scenarios:
- id: file-exists
  description: file exists
  mocks:
  - description: the file /path/to/file exists
    command: cat
    args:
    - /path/to/file
    exit-code: 1
    stdout: some text in the file`)

	singleScenarioID := "file-exists"
	Scenarios(script, scenarios, &singleScenarioID)

	if !exitCalled {
		t.Fatalf("os.Exit was not called")
	}
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
}

func TestScenariosWithExpectPassing(t *testing.T) {
	script := strings.NewReader(`>&2 echo -n "error"; cat /path/to/file`)
	scenarios := strings.NewReader(`
scenarios:
- id: file-exists
  description: file exists
  mocks:
  - description: the file /path/to/file exists
    command: cat
    args:
    - /path/to/file
    exit-code: 0
    stdout: some text in the file
  expect:
    exit-code: 0
    stdout: some text in the file
    stderr: "error"`)

	singleScenarioID := "file-exists"

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

	Scenarios(script, scenarios, &singleScenarioID)

	wStdout.Close()
	wStderr.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr
	outStdout := <-outCStdout
	outStderr := <-outCStderr

	if outStdout != "" {
		t.Fatalf("expected stdout to be empty, got %s", outStdout)
	}
	if outStderr != "" {
		t.Fatalf("expected stderr to be empty, got %s", outStderr)
	}
}

func TestScenariosWithCallsExpectationsPassing(t *testing.T) {
	script := strings.NewReader(`cat /path/to/file`)
	scenarios := strings.NewReader(`
scenarios:
- id: cat-is-called
  description: cat is called
  expect:
    calls:
    - command: cat
      args:
      - /path/to/file
      called: 1`)

	singleScenarioID := "cat-is-called"
	Scenarios(script, scenarios, &singleScenarioID)
}
