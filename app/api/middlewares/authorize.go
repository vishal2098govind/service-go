package middlewares

import (
	"context"

	"github.com/vishal2098govind/service/app/api/auth"
	"github.com/vishal2098govind/service/app/api/errs"
	"github.com/vishal2098govind/service/foundations/logger"
)

func Authorize(ctx context.Context, log *logger.Logger, auth *auth.Auth, rule string, handler Handler) error {

	claims, err := GetClaims(ctx)
	if err != nil {
		log.Error(ctx, "claims not found in context")
		return errs.Newf(errs.InvalidArgument, "claims not found: %s", err)
	}

	userId, err := GetUserID(ctx)
	if err != nil {
		log.Error(ctx, "userID not found in context")
		return errs.Newf(errs.InvalidArgument, "userID not found: %s", err)
	}

	if authorize, err := auth.Authorize(ctx, claims, rule, userId); err != nil || !authorize {
		log.Error(ctx, "unauthorized", "authorize: %s", err)
		return errs.Newf(errs.PermissionDenied, "unauthorized request")
	}

	return handler(ctx)
}
