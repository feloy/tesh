package call

import (
	"context"
	"testing"

	"github.com/feloy/tesh/pkg/scenarios"
)

func TestGetCallHandlerWithOneCall(t *testing.T) {
	calls := []scenarios.Call{
		{Command: "cat", Args: []string{"/path/to/file"}, Called: 1},
	}
	callHandler, callsResult := GetCallHandler(calls)
	callHandler(context.TODO(), []string{"cat", "/path/to/file"})
	callsResult.CheckResults()
}

func TestGetCallHandlerWithTwoCalls(t *testing.T) {
	calls := []scenarios.Call{
		{Command: "cat", Args: []string{"/path/to/file"}, Called: 2},
	}
	callHandler, callsResult := GetCallHandler(calls)
	callHandler(context.TODO(), []string{"cat", "/path/to/file"})
	callHandler(context.TODO(), []string{"cat", "/path/to/file"})
	callsResult.CheckResults()
}

func TestGetCallHandlerWithTwoCallsForTwoCommands(t *testing.T) {
	calls := []scenarios.Call{
		{Command: "cat", Args: []string{"/path/to/file"}, Called: 2},
		{Command: "ls", Args: []string{"/path/to/file"}, Called: 2},
	}
	callHandler, callsResult := GetCallHandler(calls)
	callHandler(context.TODO(), []string{"cat", "/path/to/file"})
	callHandler(context.TODO(), []string{"cat", "/path/to/file"})
	callHandler(context.TODO(), []string{"ls", "/path/to/file"})
	callHandler(context.TODO(), []string{"ls", "/path/to/file"})
	callsResult.CheckResults()
}

func TestGetCallHandlerWithNoCall(t *testing.T) {
	calls := []scenarios.Call{
		{Command: "cat", Args: []string{"/path/to/file"}, Called: 0},
	}
	_, callsResult := GetCallHandler(calls)
	callsResult.CheckResults()
}
