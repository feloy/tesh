package call

import (
	"context"
	"log"
	"reflect"

	"github.com/feloy/tesh/pkg/scenarios"
	"mvdan.cc/sh/v3/interp"
)

type CallsResult interface {
	CheckResults()
}

func GetCallHandler(calls []scenarios.Call) (interp.CallHandlerFunc, CallsResult) {
	callsResult := &callsResultImpl{
		expectedCalls: calls,
	}

	handler := func(ctx context.Context, args []string) ([]string, error) {
		for _, call := range calls {
			if call.Command == args[0] && reflect.DeepEqual(call.Args, args[1:]) {
				callsResult.addMatch(call)
				// do not break here, because we want to match all expected calls, not only the first one
			}
		}
		return args, nil
	}
	return handler, callsResult
}

type callsResultImpl struct {
	expectedCalls []scenarios.Call
	realCalls     []scenarios.Call
}

func (c *callsResultImpl) CheckResults() {
	remainingExpectedCalls := []scenarios.Call{}
	for _, expectedCall := range c.expectedCalls {
		found := false
		for _, realCall := range c.realCalls {
			if expectedCall.Command == realCall.Command && reflect.DeepEqual(expectedCall.Args, realCall.Args) {
				found = true
				if expectedCall.Called != realCall.Called {
					log.Fatalf("expected call %s %v to be called %d times, got %d times", expectedCall.Command, expectedCall.Args, expectedCall.Called, realCall.Called)
				}
			}
		}
		if !found {
			remainingExpectedCalls = append(remainingExpectedCalls, expectedCall)
		}
	}
	if len(remainingExpectedCalls) > 0 {
		for _, remainingCall := range remainingExpectedCalls {
			if remainingCall.Called != 0 {
				log.Fatalf("expected call %s %v to be called %d times, got %d times", remainingCall.Command, remainingCall.Args, remainingCall.Called, 0)
			}
		}
	}
}

func (c *callsResultImpl) addMatch(call scenarios.Call) {
	for i := range c.realCalls {
		if c.realCalls[i].Command == call.Command && reflect.DeepEqual(c.realCalls[i].Args, call.Args) {
			c.realCalls[i].Called++
			return
		}
	}
	c.realCalls = append(c.realCalls, scenarios.Call{
		Command: call.Command,
		Args:    call.Args,
		Called:  1,
	})
}
