package middlewares

import (
	"context"

	"github.com/google/uuid"
	"github.com/vishal2098govind/service/app/api/auth"
	"github.com/vishal2098govind/service/app/api/errs"
	"github.com/vishal2098govind/service/foundations/logger"
)

func Bearer(ctx context.Context, log *logger.Logger, auth *auth.Auth, authorization string, handler Handler) error {

	claims, err := auth.Authenticate(authorization)
	if err != nil {
		log.Error(ctx, "failed to authenticate token", "authorization", authorization, "err", err)
		return errs.Newf(errs.Unauthenticated, "authenticate: request is not authenticated")
	}

	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		return errs.Newf(errs.Unauthenticated, "invalid subject. must be a valid uuid")
	}
	log.Info(ctx, "authenticated request", "userId", userId)

	ctx = setUserID(ctx, userId)
	ctx = setClaims(ctx, claims)

	return handler(ctx)
}

func Basic(ctx context.Context, log *logger.Logger, ath *auth.Auth, username string, pass string, handler Handler) error {

	claims, err := ath.Basic(username, pass)
	if err != nil {
		if err == auth.ErrInvalidCreds {
			return errs.Newf(errs.InvalidArgument, "invalid credentials")
		}
		return errs.New(errs.Internal, "failed to authenticate user")
	}

	ctx = setClaims(ctx, claims)

	return handler(ctx)
}
