package output

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"mvdan.cc/sh/v3/syntax"
)

func TestOutputCoverage(t *testing.T) {
	coverage := []syntax.Pos{
		syntax.NewPos(0, 1, 1),
		syntax.NewPos(10, 1, 11),
		syntax.NewPos(30, 2, 4),
		syntax.NewPos(41, 2, 15),
	}
	lens := []uint{6, 6, 5, 7}
	covered := []uint{1, 1, 1, 1}
	file := strings.NewReader(`echo A && echo B || echo C
if true; then echo D; else echo E; fi`)
	buf := bytes.Buffer{}
	OutputCoverage(&buf, file, coverage, lens, covered)
	expected := fmt.Sprintf(`%secho A%s && %secho B%s || echo C
if %strue;%s then %secho D;%s else echo E; fi`, COVERED_COLOR, RESET_COLOR, COVERED_COLOR, RESET_COLOR, COVERED_COLOR, RESET_COLOR, COVERED_COLOR, RESET_COLOR)
	if buf.String() != expected {
		t.Fatalf("expected stdout to be\n%s, got\n%s", expected, buf.String())
	}
}
