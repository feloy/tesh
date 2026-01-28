package expect

import (
	"reflect"
	"strings"
	"testing"

	"github.com/feloy/tesh/pkg/api"
	"github.com/feloy/tesh/pkg/scenarios"
)

func TestCheckExpectations(t *testing.T) {
	tests := []struct {
		name           string
		expectations   scenarios.Expect
		exitCode       int
		stdout         string
		stderr         string
		expectedResult api.ScenarioResult
	}{
		{
			name: "exit code OK",
			expectations: scenarios.Expect{
				ExitCode: &[]int{0}[0],
			},
			exitCode:       0,
			stdout:         "",
			stderr:         "",
			expectedResult: api.ScenarioResult{},
		}, {
			name: "exit code wrong",
			expectations: scenarios.Expect{
				ExitCode: &[]int{0}[0],
			},
			exitCode: 1,
			stdout:   "",
			stderr:   "",
			expectedResult: api.ScenarioResult{
				ExitCode: &api.ExitCodeError{
					Expected: 0,
					Actual:   1,
				},
			},
		}, {
			name: "non empty stdout OK",
			expectations: scenarios.Expect{
				Stdout: &[]string{"some text"}[0],
			},
			exitCode:       0,
			stdout:         "some text",
			stderr:         "",
			expectedResult: api.ScenarioResult{},
		}, {
			name: "empty stdout OK",
			expectations: scenarios.Expect{
				Stdout: &[]string{""}[0],
			},
			exitCode:       0,
			stdout:         "",
			stderr:         "",
			expectedResult: api.ScenarioResult{},
		}, {
			name: "non empty stdout wrong",
			expectations: scenarios.Expect{
				Stdout: &[]string{"some text in standard output"}[0],
			},
			exitCode: 0,
			stdout:   "other text in standard output",
			stderr:   "",
			expectedResult: api.ScenarioResult{
				Stdout: &api.StdioError{
					Expected: "some text in standard output",
					Actual:   "other text in standard output",
				},
			},
		}, {
			name: "empty stdout wrong",
			expectations: scenarios.Expect{
				Stdout: &[]string{""}[0],
			},
			exitCode: 0,
			stdout:   "other text in standard output",
			stderr:   "",
			expectedResult: api.ScenarioResult{
				Stdout: &api.StdioError{
					Expected: "",
					Actual:   "other text in standard output",
				},
			},
		}, {
			name: "non empty stdout wrong 2",
			expectations: scenarios.Expect{
				Stdout: &[]string{"some text"}[0],
			},
			exitCode: 0,
			stdout:   "",
			stderr:   "",
			expectedResult: api.ScenarioResult{
				Stdout: &api.StdioError{
					Expected: "some text",
					Actual:   "",
				},
			},
		}, {
			name: "non empty stderr OK",
			expectations: scenarios.Expect{
				Stderr: &[]string{"some text"}[0],
			},
			exitCode:       0,
			stdout:         "",
			stderr:         "some text",
			expectedResult: api.ScenarioResult{},
		}, {
			name: "empty stdout OK",
			expectations: scenarios.Expect{
				Stderr: &[]string{""}[0],
			},
			exitCode:       0,
			stdout:         "",
			stderr:         "",
			expectedResult: api.ScenarioResult{},
		}, {
			name: "non empty stderr wrong",
			expectations: scenarios.Expect{
				Stderr: &[]string{"some text in standard error"}[0],
			},
			exitCode: 0,
			stdout:   "",
			stderr:   "other text in standard error",
			expectedResult: api.ScenarioResult{
				Stderr: &api.StdioError{
					Expected: "some text in standard error",
					Actual:   "other text in standard error",
				},
			},
		}, {
			name: "empty stderr wrong",
			expectations: scenarios.Expect{
				Stderr: &[]string{""}[0],
			},
			exitCode: 0,
			stdout:   "",
			stderr:   "other text in standard error",
			expectedResult: api.ScenarioResult{
				Stderr: &api.StdioError{
					Expected: "",
					Actual:   "other text in standard error",
				},
			},
		}, {
			name: "non empty stderr wrong 2",
			expectations: scenarios.Expect{
				Stderr: &[]string{"some text"}[0],
			},
			exitCode: 0,
			stdout:   "",
			stderr:   "",
			expectedResult: api.ScenarioResult{
				Stderr: &api.StdioError{
					Expected: "some text",
					Actual:   "",
				},
			},
		}, {
			name: "everything OK",
			expectations: scenarios.Expect{
				ExitCode: &[]int{0}[0],
				Stdout:   &[]string{"some stdout"}[0],
				Stderr:   &[]string{"some stderr"}[0],
			},
			exitCode:       0,
			stdout:         "some stdout",
			stderr:         "some stderr",
			expectedResult: api.ScenarioResult{},
		}, {
			name: "everything wrong",
			expectations: scenarios.Expect{
				ExitCode: &[]int{0}[0],
				Stdout:   &[]string{"some stdout"}[0],
				Stderr:   &[]string{"some stderr"}[0],
			},
			exitCode: 1,
			stdout:   "other stdout",
			stderr:   "other stderr",
			expectedResult: api.ScenarioResult{
				ExitCode: &api.ExitCodeError{
					Expected: 0,
					Actual:   1,
				},
				Stdout: &api.StdioError{
					Expected: "some stdout",
					Actual:   "other stdout",
				},
				Stderr: &api.StdioError{
					Expected: "some stderr",
					Actual:   "other stderr",
				},
			},
		},
	}
	for _, test := range tests {
		result := &api.ScenarioResult{}
		CheckExpectations(&test.expectations, result, test.exitCode, strings.NewReader(test.stdout), strings.NewReader(test.stderr))
		if !reflect.DeepEqual(*result, test.expectedResult) {
			t.Fatalf("%s: expected result to be %v, got %v", test.name, test.expectedResult, result)
		}
	}
}
