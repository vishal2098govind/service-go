package web

import (
	"context"
	"time"
)

type ctxKey int

const key ctxKey = 1

type Value struct {
	TraceID string
	Now     time.Time
}

func GetValue(ctx context.Context) *Value {
	v, ok := ctx.Value(key).(*Value)
	if !ok {
		return &Value{
			TraceID: "00000000-0000-0000-0000-000000000000",
			Now:     time.Now(),
		}
	}

	return v
}

func GetTraceID(ctx context.Context) string {
	v := GetValue(ctx)
	return v.TraceID
}

func setValues(ctx context.Context, value *Value) context.Context {
	ctx = context.WithValue(ctx, key, value)
	return ctx
}
