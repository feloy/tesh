package api

type ScenarioResult struct {
	ScenarioID string
	ExitCode   *ExitCodeError
	Stdout     *StdioError
	Stderr     *StdioError
	Calls      []CallResult
}

type ExitCodeError struct {
	Expected int
	Actual   int
}

type StdioError struct {
	Expected string
	Actual   string
}

type CallResult struct {
	Command        string
	Args           []string
	ExpectedCalled int
	ActualCalled   int
}

func (o *ScenarioResult) IsSuccess() bool {
	return o.ExitCode == nil && o.Stdout == nil && o.Stderr == nil && len(o.Calls) == 0
}
