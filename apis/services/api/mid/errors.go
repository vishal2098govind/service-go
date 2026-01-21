package mid

import (
	"context"
	"net/http"

	"github.com/vishal2098govind/service/app/api/mid"
	"github.com/vishal2098govind/service/foundations/logger"
	"github.com/vishal2098govind/service/foundations/web"
)

// API layer Errors middleware
func Errors(log *logger.Logger) web.MidHandler {
	return func(handler web.Handler) web.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			hld := func(ctx context.Context) error {
				return handler(ctx, w, r)
			}

			err := mid.Errors(ctx, log, hld)
			if err != nil {
				if err := web.Respond(ctx, w, struct {
					Code string
				}{
					Code: http.StatusText(http.StatusInternalServerError),
				}, http.StatusInternalServerError); err != nil {
					return err
				}
			}

			return nil
		}
	}
}
