package middlewares

import (
	"context"

	"github.com/vishal2098govind/service/app/api/errs"
	"github.com/vishal2098govind/service/foundations/logger"
)

// App layer Error middleware
func Errors(ctx context.Context, log *logger.Logger, handler Handler) error {
	err := handler(ctx)

	if err == nil {
		return nil
	}

	log.Error(ctx, "ERROR", "error", err.Error())

	if errs.IsError(err) {
		return errs.GetError(err)
	}

	return errs.New(errs.Unknown, errs.Unknown.String())
}
