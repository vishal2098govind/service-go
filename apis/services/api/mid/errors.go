package mid

import (
	"context"
	"net/http"

	"github.com/vishal2098govind/service/app/api/errs"
	mid "github.com/vishal2098govind/service/app/api/middlewares"
	"github.com/vishal2098govind/service/foundations/logger"
	"github.com/vishal2098govind/service/foundations/web"
)

var codeStatuses [17]int

func init() {
	codeStatuses[errs.OK.Value()] = http.StatusOK
	codeStatuses[errs.Canceled.Value()] = http.StatusGatewayTimeout
	codeStatuses[errs.Unknown.Value()] = http.StatusInternalServerError
	codeStatuses[errs.InvalidArgument.Value()] = http.StatusBadRequest
	codeStatuses[errs.DeadlineExceeded.Value()] = http.StatusGatewayTimeout
	codeStatuses[errs.NotFound.Value()] = http.StatusNotFound
	codeStatuses[errs.AlreadyExists.Value()] = http.StatusConflict
	codeStatuses[errs.PermissionDenied.Value()] = http.StatusForbidden
	codeStatuses[errs.ResourceExhausted.Value()] = http.StatusTooManyRequests
	codeStatuses[errs.FailedPrecondition.Value()] = http.StatusBadRequest
	codeStatuses[errs.Aborted.Value()] = http.StatusConflict
	codeStatuses[errs.OutOfRange.Value()] = http.StatusBadRequest
	codeStatuses[errs.Unimplemented.Value()] = http.StatusNotImplemented
	codeStatuses[errs.Internal.Value()] = http.StatusInternalServerError
	codeStatuses[errs.Unavailable.Value()] = http.StatusServiceUnavailable
	codeStatuses[errs.DataLoss.Value()] = http.StatusInternalServerError
	codeStatuses[errs.Unauthenticated.Value()] = http.StatusUnauthorized
}

// API layer Errors middleware
func Errors(log *logger.Logger) web.MidHandler {
	return func(handler web.Handler) web.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			hld := func(ctx context.Context) error {
				return handler(ctx, w, r)
			}

			err := mid.Errors(ctx, log, hld)
			if err != nil {
				er := err.(errs.Error)
				if err := web.Respond(ctx, w, err, codeStatuses[er.Code.Value()]); err != nil {
					return err
				}
			}

			return nil
		}
	}
}
