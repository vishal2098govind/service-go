package middlewares

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/vishal2098govind/service/foundations/logger"
)

func Panics(ctx context.Context, log *logger.Logger, hld Handler) (err error) {

	defer func() {
		if rec := recover(); rec != nil {
			trace := debug.Stack()
			err = fmt.Errorf("PANIC [%v] STACK: %s", rec, string(trace))
		}
	}()

	return hld(ctx)
}
