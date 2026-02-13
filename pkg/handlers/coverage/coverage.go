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
	ch         chan syntax.Pos
	done       chan struct{}
}

func New(script *syntax.File) *Coverage {
	o := &Coverage{
		script:     script,
		statements: make(map[syntax.Pos]syntax.Node),
		covered:    make(map[syntax.Pos]bool),
		ch:         make(chan syntax.Pos),
		done:       make(chan struct{}),
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
	go func() {
		for {
			pos, more := <-o.ch
			if !more {
				o.done <- struct{}{}
				return
			}
			o.covered[pos] = true
		}
	}()
}

func (o *Coverage) setCovered(pos syntax.Pos) {
	// avoid conflict by using a Go channel
	o.ch <- pos
}

func (o *Coverage) GetCoverageHandler() interp.CallHandlerFunc {
	handler := func(ctx context.Context, args []string) ([]string, error) {
		hc := interp.HandlerCtx(ctx)
		o.setCovered(hc.Pos)
		return args, nil
	}
	return handler
}

func (o *Coverage) GetCoverageResult() ([]syntax.Pos, []uint, []uint) {
	close(o.ch)
	<-o.done

	keys := make([]syntax.Pos, 0, len(o.statements))
	for pos := range o.statements {
		keys = append(keys, pos)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Offset() < keys[j].Offset()
	})

	lens := make([]uint, 0, len(o.statements))
	covered := make([]uint, 0, len(o.statements))
	for _, pos := range keys {
		lens = append(lens, o.statements[pos].End().Offset()-pos.Offset())
		if o.covered[pos] {
			covered = append(covered, 1)
		} else {
			covered = append(covered, 0)
		}
	}
	return keys, lens, covered
}
