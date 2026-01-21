package mid

import (
	"context"
	"net/http"

	"github.com/vishal2098govind/service/app/api/mid"
	"github.com/vishal2098govind/service/foundations/logger"
	"github.com/vishal2098govind/service/foundations/web"
)

func Panics(log *logger.Logger) web.MidHandler {
	return func(handler web.Handler) web.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			hld := func(ctx context.Context) error {
				return handler(ctx, w, r)
			}

			return mid.Panics(ctx, log, hld)
		}
	}
}
