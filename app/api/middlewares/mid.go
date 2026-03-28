package middlewares

import (
	"context"
)

type Handler func(context.Context) error
