package coverage

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

type TestContext struct {
	pos syntax.Pos
}

func (t *TestContext) Value(key any) any {
	return interp.HandlerContext{
		Pos: t.pos,
	}
}

func (t *TestContext) Deadline() (time.Time, bool) {
	return time.Time{}, false
}

func (t *TestContext) Done() <-chan struct{} {
	return nil
}

func (t *TestContext) Err() error {
	return nil
}

func TestCoverage(t *testing.T) {
	file := strings.NewReader(`echo A && echo B || echo C
if true; then echo D; else echo E; fi`)
	script, _ := syntax.NewParser().Parse(file, "")
	coverage := New(script)
	handler := coverage.GetCoverageHandler()

	ctx := &TestContext{
		pos: syntax.NewPos(0, 1, 1),
	}
	handler(ctx, []string{"echo", "A"})

	ctx = &TestContext{
		pos: syntax.NewPos(10, 1, 11),
	}
	handler(ctx, []string{"echo", "B"})

	ctx = &TestContext{
		pos: syntax.NewPos(30, 2, 4),
	}
	handler(ctx, []string{"true"})

	ctx = &TestContext{
		pos: syntax.NewPos(41, 2, 15),
	}
	handler(ctx, []string{"echo", "D"})

	positions, lens, covered := coverage.GetCoverageResult()
	if !reflect.DeepEqual(positions, []syntax.Pos{
		syntax.NewPos(0, 1, 1),
		syntax.NewPos(10, 1, 11),
		syntax.NewPos(20, 1, 21),
		syntax.NewPos(30, 2, 4),
		syntax.NewPos(41, 2, 15),
		syntax.NewPos(54, 2, 28),
	}) {
		t.Fatalf("incorrect expected positions, got %v", positions)
	}
	if !reflect.DeepEqual(lens, []uint{6, 6, 6, 5, 7, 7}) {
		t.Fatalf("incorrect expected lens, got %v", lens)
	}
	if !reflect.DeepEqual(covered, []uint{1, 1, 0, 1, 1, 0}) {
		t.Fatalf("incorrect expected covered, got %v", covered)
	}
}
