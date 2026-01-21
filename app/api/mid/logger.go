package mid

import (
	"context"
	"fmt"

	"github.com/vishal2098govind/service/foundations/logger"
)

type Handler func(context.Context) error

// this logger middleware is protocol agnostic
func Logger(ctx context.Context, log *logger.Logger, path string, rawQuery string, method string, remoteAddress string, handler Handler) error {

	if rawQuery != "" {
		path = fmt.Sprintf("%s?%s", path, rawQuery)
	}

	log.Info(ctx, "request started", "path", path, "method", method, "remoteAddress", remoteAddress)

	err := handler(ctx)

	log.Info(ctx, "request completed", "path", path, "method", method, "remoteAddress", remoteAddress)

	return err

}
