package exec

import (
	"context"
	"fmt"

	"github.com/feloy/tesh/pkg/scenarios"
	"mvdan.cc/sh/v3/interp"
)

func GetExecHandler(mock scenarios.Mock) func(next interp.ExecHandlerFunc) interp.ExecHandlerFunc {
	return func(next interp.ExecHandlerFunc) interp.ExecHandlerFunc {
		return func(ctx context.Context, args []string) error {
			if args[0] != mock.Command {
				return next(ctx, args)
			}
			for i, arg := range mock.Args {
				if arg != args[i+1] {
					return next(ctx, args)
				}
			}
			if mock.Stdout != nil {
				hc := interp.HandlerCtx(ctx)
				fmt.Fprintf(hc.Stdout, "%s", *mock.Stdout)
			}
			if mock.Stderr != nil {
				hc := interp.HandlerCtx(ctx)
				fmt.Fprintf(hc.Stderr, "%s", *mock.Stderr)
			}
			if mock.ExitCode != nil {
				return interp.ExitStatus(*mock.ExitCode)
			} else {
				return nil
			}
		}
	}
}
