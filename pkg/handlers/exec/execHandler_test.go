package exec

import (
	"context"
	"errors"
	"testing"

	"github.com/feloy/tesh/pkg/scenarios"
	"mvdan.cc/sh/v3/interp"
)

func TestExitCode1(t *testing.T) {
	mock := scenarios.Mock{
		Description: "git diff indicates there are changes in the repository",
		Command:     "git",
		Args:        []string{"diff"},
		ExitCode:    &[]int{1}[0],
	}
	execHandler := GetExecHandler(mock)
	f := execHandler(func(ctx context.Context, args []string) error {
		return nil
	})
	err := f(context.TODO(), []string{"git", "diff"})
	if err != interp.ExitStatus(1) {
		t.Fatalf("expected exit status 1, got %v", err)
	}
}

func TestExitCode0(t *testing.T) {
	mock := scenarios.Mock{
		Description: "git diff indicates there are changes in the repository",
		Command:     "git",
		Args:        []string{"diff"},
		ExitCode:    &[]int{0}[0],
	}
	execHandler := GetExecHandler(mock)
	f := execHandler(func(ctx context.Context, args []string) error {
		return nil
	})
	err := f(context.TODO(), []string{"git", "diff"})
	if err != interp.ExitStatus(0) {
		t.Fatalf("expected exit status 0, got %v", err)
	}
}

func TestExitCodeNil(t *testing.T) {
	mock := scenarios.Mock{
		Description: "git diff indicates there are changes in the repository",
		Command:     "git",
		Args:        []string{"diff"},
	}
	execHandler := GetExecHandler(mock)
	f := execHandler(func(ctx context.Context, args []string) error {
		return nil
	})
	err := f(context.TODO(), []string{"git", "diff"})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestOtherCommand(t *testing.T) {
	mock := scenarios.Mock{
		Description: "git diff indicates there are changes in the repository",
		Command:     "git",
		Args:        []string{"diff"},
		ExitCode:    &[]int{1}[0],
	}
	execHandler := GetExecHandler(mock)
	nextError := errors.New("the error from next handler")
	next := func(ctx context.Context, args []string) error {
		return nextError
	}
	f := execHandler(next)
	err := f(context.TODO(), []string{"ls"})
	if err != nextError {
		t.Fatalf("expected nextError, got %v", err)
	}
}

func TestOtherArgs(t *testing.T) {
	mock := scenarios.Mock{
		Description: "git diff indicates there are changes in the repository",
		Command:     "git",
		Args:        []string{"diff"},
		ExitCode:    &[]int{1}[0],
	}
	execHandler := GetExecHandler(mock)
	nextError := errors.New("the error from next handler")
	next := func(ctx context.Context, args []string) error {
		return nextError
	}
	f := execHandler(next)
	err := f(context.TODO(), []string{"git", "status"})
	if err != nextError {
		t.Fatalf("expected nextError, got %v", err)
	}
}
