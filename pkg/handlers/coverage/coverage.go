package coverage

import (
	"context"
	"sort"

	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

type Coverage struct {
	script     *syntax.File
	statements map[syntax.Pos]syntax.Node
	covered    map[syntax.Pos]bool
}

func New(script *syntax.File) *Coverage {
	o := &Coverage{
		script:     script,
		statements: make(map[syntax.Pos]syntax.Node),
		covered:    make(map[syntax.Pos]bool),
	}
	o.init()
	return o
}

func (o *Coverage) init() {
	syntax.Walk(o.script, func(node syntax.Node) bool {
		switch node.(type) {
		case *syntax.Stmt:
			o.statements[node.Pos()] = node
		}
		return true
	})
}

func (o *Coverage) GetCoverageHandler() interp.CallHandlerFunc {
	handler := func(ctx context.Context, args []string) ([]string, error) {
		hc := interp.HandlerCtx(ctx)
		o.covered[hc.Pos] = true
		return args, nil
	}
	return handler
}

func (o *Coverage) GetCoverageResult() ([]syntax.Pos, []uint) {
	keys := make([]syntax.Pos, 0, len(o.covered))
	for pos := range o.covered {
		keys = append(keys, pos)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Offset() < keys[j].Offset()
	})
	lens := make([]uint, 0, len(o.covered))
	for _, pos := range keys {
		lens = append(lens, o.statements[pos].End().Offset()-pos.Offset())
	}
	return keys, lens
}
