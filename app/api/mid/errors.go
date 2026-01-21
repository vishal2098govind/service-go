package mid

import (
	"context"

	"github.com/vishal2098govind/service/foundations/logger"
)

// App layer Error middleware
func Errors(ctx context.Context, log *logger.Logger, handler Handler) error {
	err := handler(ctx)

	if err == nil {
		return nil
	}

	log.Error(ctx, "ERROR", "error", err.Error())

	return err
}
