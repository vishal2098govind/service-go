package middlewares

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/vishal2098govind/service/app/api/auth"
)

type Handler func(context.Context) error

type ctxKey int

const (
	userCtx ctxKey = iota + 1
	claimsCtx
)

func setUserID(ctx context.Context, id uuid.UUID) context.Context {
	ctx = context.WithValue(ctx, userCtx, id)
	return ctx
}

func GetUserID(ctx context.Context) (uuid.UUID, error) {
	userIdVal, ok := ctx.Value(userCtx).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("user id not found in context")
	}
	return userIdVal, nil
}

func setClaims(ctx context.Context, claims auth.Claims) context.Context {
	return context.WithValue(ctx, claimsCtx, claims)
}

func GetClaims(ctx context.Context) (auth.Claims, error) {
	claims, ok := ctx.Value(claimsCtx).(auth.Claims)
	if !ok {
		return auth.Claims{}, fmt.Errorf("claims not found in context")
	}
	return claims, nil
}
