package middlewares

import (
	"context"
	"net/http"

	"github.com/vishal2098govind/service/app/api/auth"
	"github.com/vishal2098govind/service/app/api/middlewares"
	"github.com/vishal2098govind/service/foundations/logger"
	"github.com/vishal2098govind/service/foundations/web"
)

func Authenticate(log *logger.Logger, auth *auth.Auth) web.MidHandler {
	return func(h web.Handler) web.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			hdl := func(ctx context.Context) error {
				return h(ctx, w, r)
			}
			authorization := r.Header.Get("authorization")
			err := middlewares.Authenticate(ctx, log, auth, authorization, hdl)
			return err
		}
	}
}
