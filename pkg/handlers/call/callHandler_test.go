package call

import (
	"context"
	"testing"

	"github.com/feloy/tesh/pkg/api"
	"github.com/feloy/tesh/pkg/scenarios"
)

func TestGetCallHandler(t *testing.T) {
	tests := []struct {
		name                   string
		expectedCalls          []scenarios.Call
		realCalls              [][]string
		expectedScenarioResult api.ScenarioResult
	}{
		{
			name: "one matching call",
			expectedCalls: []scenarios.Call{
				{Command: "cat", Args: []string{"/path/to/file"}, Called: 1},
			},
			realCalls: [][]string{
				{"cat", "/path/to/file"},
			},
			expectedScenarioResult: api.ScenarioResult{},
		},
		{
			name: "two matching calls",
			expectedCalls: []scenarios.Call{
				{Command: "cat", Args: []string{"/path/to/file"}, Called: 2},
			},
			realCalls: [][]string{
				{"cat", "/path/to/file"},
				{"cat", "/path/to/file"},
			},
			expectedScenarioResult: api.ScenarioResult{},
		},
		{
			name: "two different matching calls",
			expectedCalls: []scenarios.Call{
				{Command: "cat", Args: []string{"/path/to/file"}, Called: 2},
				{Command: "ls", Args: []string{"/path/to/other/file"}, Called: 1},
			},
			realCalls: [][]string{
				{"cat", "/path/to/file"},
				{"cat", "/path/to/file"},
				{"ls", "/path/to/other/file"},
			},
			expectedScenarioResult: api.ScenarioResult{},
		},
		{
			name: "expect no call",
			expectedCalls: []scenarios.Call{
				{Command: "cat", Args: []string{"/path/to/file"}, Called: 0},
			},
			realCalls: [][]string{
				{"ls", "/path/to/other/file"},
			},
			expectedScenarioResult: api.ScenarioResult{},
		},
		{
			name: "one non-matching call",
			expectedCalls: []scenarios.Call{
				{Command: "cat", Args: []string{"/path/to/file"}, Called: 1},
			},
			realCalls: [][]string{
				{"cat", "/path/to/other/file"},
			},
			expectedScenarioResult: api.ScenarioResult{
				Calls: []api.CallResult{
					{Command: "cat", Args: []string{"/path/to/file"}, ExpectedCalled: 1, ActualCalled: 0},
				},
			},
		},
		{
			name: "unexpected call",
			expectedCalls: []scenarios.Call{
				{Command: "cat", Args: []string{"/path/to/file"}, Called: 0},
			},
			realCalls: [][]string{
				{"cat", "/path/to/file"},
			},
			expectedScenarioResult: api.ScenarioResult{
				Calls: []api.CallResult{
					{Command: "cat", Args: []string{"/path/to/file"}, ExpectedCalled: 0, ActualCalled: 1},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			callHandler, callsResult := GetCallHandler(test.expectedCalls)
			for _, realCall := range test.realCalls {
				callHandler(context.TODO(), realCall)
			}
			callsResult.CheckResults(&test.expectedScenarioResult)
		})
	}
}
