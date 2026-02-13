# tesh

tl;dr: Testing shell scripts made easy. Run scripts with mocks, assertions, and coverage.

![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)
[![codecov](https://codecov.io/gh/feloy/tesh/graph/badge.svg)](https://codecov.io/gh/feloy/tesh)

## Features

- Shell interpreter
- Mock external commands with specific outputs (stdout and stderr) and exit codes
- Mock environment variables
- Mock file existence
- Assert script's exit code and output (stdout and stderr)
- Assert calls of external commands
- Generate colored coverage output or Go coverage profiles.

## Install

```sh
go install github.com/feloy/tesh@latest
```

Or build from source:

```sh
go build -o tesh .
```

## Shell Interpreter

The `tesh` command relies on the mvdan's Go `sh` library (https://github.com/mvdan/sh) to interpret and run shell scripts. `tesh` uses the many handlers provided by this interpreter to implement mocks, assertions and coverage.

The limitations of the interpreter are described in the library documentation (https://pkg.go.dev/mvdan.cc/sh/v3/interp):

> Package interp implements an interpreter that executes shell programs. It aims to support POSIX, but its support is not complete yet. It also supports some Bash features.

> The interpreter generally aims to behave like Bash, but it does not support all of its features.

> The interpreter currently aims to behave like a non-interactive shell, which is how most shells run scripts, and is more useful to machines. In the future, it may gain an option to behave like an interactive shell.

To be sure that your scripts are executed in the exact same way they are tested, you can run the scripts "in production" with the `tesh` command:

Script (`script.sh`):
```sh
echo Hello World
```

Console:
```console
$ tesh script.sh
Hello World
```

## Mocking external commands

You can provide a _Scenario_ file to `tesh` with the `--scenarios` flag. This file defines one or more scenarios, and each scenario can define mocks for external commands among other things.

When scenarios only define mocks and no expectations, `tesh` behaves as a normal shell interpreter, except that instead of executing mocked sub-commands, it uses the exit code and stdout and stderr defined in the scenario.

In this mode, if several scenarios are defined in the Scenario file, and you do not scpecify a scenario ID, the interpreter executes the script with the first scenario; or you can specify the scenario to execute with the `--scenario` flag.


Script (`examples/ex1.sh`):
```sh
cat /path/to/file
```

Scenarios (`examples/ex1.yaml`):
```yaml
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

- id: file-not-exists
  description: file does not exist
  mocks:
  - description: the file /path/to/file does not exist
    command: cat
    args:
    - /path/to/file
    exit-code: 1
    stderr: the file /path/to/file does not exist
```

In this example, the script is executed with the first scenario, and because the sub-command `cat /path/to/file` is mocked, the interpreter does not call the `cat` command, but behaves as if the command had exited with the status code 0 and had written `some text in the file` in its stdout.

Console:

```console
$ tesh examples/ex1.sh \
    --scenarios examples/ex1.yaml
some text in the file
$ echo $?
0
```

In this example, the script is executed with the scenario `file-exists`, as in the previous example.

Console:
```
$ tesh examples/ex1.sh \
    --scenarios examples/ex1.yaml \
    --scenario file-exists
some text in the file
$ echo $?
0
```

In this example, the script is executed with the scenario `file-not-exists`. The sub-command `cat /path/to/file` being mocked, the interpreter does not call the `cat` command, but behaves as if the command had exited with the status code 1 and had written `the file /path/to/file does not exist` in its stderr. Because this sub-command is the latest one in the script and terminates with a status code 1, the script terminates with the status code 1.


Console:
```
$ tesh examples/ex1.sh \
    --scenarios examples/ex1.yaml \
    --scenario file-not-exists
the file /path/to/file does not exist
$ echo $?
1
```

## Mocking Environment Variables

A scenario can define one or more environment variables.

Script (`examples/ex10.sh`):
```sh
echo -n $MYVAR
```

Scenarios (`examples/ex10.yaml`):
```yaml
scenarios:
- id: env-not-set
  description: MYVAR env is not set
- id: env-is-set
  description: MYVAR is set with myvalue
  envs:
  - MYVAR=myvalue
```

In this example, MYVAR is not defined, nothing is displayed.

Console:
```console
$ tesh examples/ex10.sh \
    --scenarios examples/ex10.yaml \
    --scenario env-not-set
```

In this example, MYVAR is mocked by the scenario:

Console:
```console
$ tesh examples/ex10.sh \
    --scenarios examples/ex10.yaml \
    --scenario env-is-set
myvalue
```

In this example, MYVAR is defined and not mocked, its original value is displayed:

Console:

```console
$ MYVAR=originalValue tesh examples/ex10.sh \
    --scenarios examples/ex10.yaml \
    --scenario env-not-set
originalValue
```

In this example, MYVAR is defined, and is also mocked by the scenario; the value displayed is the one provided by the mock:

Console:
```console
$ MYVAR=originalValue tesh examples/ex10.sh \
    --scenarios examples/ex10.yaml \
    --scenario env-is-set
myvalue
```

## Mocking file existence

A scenario can mock the existence of files.

Script (`examples/ex11.sh`):
```sh
[ -f ./path/to/file ] && echo -n "file exists" || echo -n "file does not exist"
```

Scenarios (`examples/ex11.yaml`):
```yaml
scenarios:
- id: file-exists
  description: file exists
  files:
  - path: ./path/to/file
    exists: true
- id: file-not-exists
  description: file does not exist
  files:
  - path: ./path/to/file
    exists: false
```


In this example, the file is mocked as existing:

Console:
```console
$ tesh examples/ex11.sh \
    --scenarios examples/ex11.yaml \
    --scenario file-exists
file exists
```

In this example, the file is mocked as non existing:

Console:
```console
$ tesh examples/ex11.sh \
    --scenarios examples/ex11.yaml \
    --scenario file-not-exists
file does not exist
```

## Asserting script's exit code and outputs

Scenarios can define assertions in addition to mocks. When all scenarios of a Scenario file provide assertions, the script is interpreted for each of the scenarios.

In this mode, all the scenarios are executed, and the exit code of `tesh` is 0 if and only if all assertions pass, or 1 otherwise. Also, the stdout and stderr for scenarios are discarded, and the result of the assertions are displayed instead.

In this example, the script is executed for all the scenarios of the Scenario file, and the results of the assertions are displayed:


Script (`examples/ex2.sh`):
```sh
cat /path/to/file
```

Scenarios (not complete, see complete file in [./examples/ex2.yaml](./examples/ex2.yaml)):
```yaml
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
    stderr: ""
```


Console:
```console
$ tesh examples/ex2.sh \
    --scenarios examples/ex2.yaml 
Scenario: file-exists
Scenario: file-not-exists
Scenario: file-not-exists-failing-exit-code
Exit Code: expected 0, actual 1
Scenario: file-not-exists-failing-stdout
Stdout: expected "some wrong expect", actual ""
Scenario: file-not-exists-failing-stderr
Stderr: expected "some wrong stderr", actual "the file /path/to/file does not exist"
```

## Asserting calls of external commands

A scenario can assert that a sub-command has been called a specific number of times. 


Script (`examples/ex3.sh`):
```sh
if ls /file/exists; then
    cat /file/exists
fi
```

Scenarios (not complete, see complete file in [./examples/ex3.yaml](./examples/ex3.yaml)):
```yaml
scenarios:
- id: file-exists
  description: file exists
  mocks:
  - description: the file /file/exists exists
    command: ls
    args:
    - /file/exists
    exit-code: 0
  - description: the file /file/exists has a content
    command: cat
    args:
    - /file/exists
    exit-code: 0
    stdout: some text in the file
  expect:
    exit-code: 0
    stdout: some text in the file
    stderr: ""
    calls:
    - command: cat
      args:
      - /file/exists
      called: 1
```

Console:
```console
$ tesh examples/ex3.sh \
    --scenarios examples/ex3.yaml 
Scenario: file-exists
Scenario: file-not-exists
Scenario: file-not-exists-failing-call
Call: cat [/file/exists], expected 1 calls, actual 0 calls
Scenario: file-exists-failing-call
Call: cat [/file/exists], expected 0 calls, actual 1 calls
```

## Usage

```text
tesh <script file> \
  [--scenarios <scenarios file> [--scenario <scenario id>] ] \
  [--coverage[=<file>]]
```

Flags:

- `--scenarios`: YAML file defining scenarios. If omitted, the script is only executed.
- `--scenario`: Run a single scenario by id (requires `--scenarios`).
- `--coverage`: Show colored coverage on stdout (suppresses normal stdout/stderr).
- `--coverage=<file>`: Write Go `coverage.txt` style output to the given file.

## Scenarios file format

Top-level `scenarios` is a list of test scenarios. Each scenario supports:

- `id` (string, required)
- `description` (string, optional)
- `mocks` (list) to fake external commands
- `envs` (list) to set environment variables
- `files` (list) to fake file existence
- `expect` (object) for assertions

Sub-commands Mocks:

```yaml
mocks:
  - description: the file /path/to/file exists
    command: cat
    args: [/path/to/file]
    exit-code: 0
    stdout: some text in the file
    stderr: ""
```

Environment variables:

```yaml
envs:
  - MYVAR=myvalue
```

File existence:

```yaml
files:
  - path: ./path/to/file
    exists: true
  - path: ./path/to/other/file
    exists: false
```

Assertions:

```yaml
expect:
  exit-code: 0
  stdout: "hello"
  stderr: ""
  calls:
    - command: cat
      args: [/path/to/file]
      called: 1
```

## Coverage

Use `--coverage` to render the script with covered lines highlighted in green and uncovered in red.

Use `--coverage=coverage.txt` to create a Go `coverprofile` compatible file.

Coverage can be used in any mode (interpreter mode without scenarios, or with scenarios and with or without assertions).

## License

Apache 2.0. See [LICENSE](LICENSE).
