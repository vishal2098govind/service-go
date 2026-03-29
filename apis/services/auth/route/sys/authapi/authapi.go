package authapi

import (
	"context"
	"net/http"

	"github.com/vishal2098govind/service/app/api/auth"
	"github.com/vishal2098govind/service/app/api/errs"
	"github.com/vishal2098govind/service/app/api/middlewares"
	"github.com/vishal2098govind/service/foundations/logger"
	"github.com/vishal2098govind/service/foundations/web"
)

type api struct {
	log  *logger.Logger
	auth *auth.Auth
}

func newAPI(auth *auth.Auth, log *logger.Logger) api {
	return api{auth: auth, log: log}
}

func (api *api) generateToken(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	kid := r.PathValue("kid")

	claims, err := middlewares.GetClaims(ctx)
	if err != nil {
		return errs.New(errs.Internal, "claims not found")
	}

	token, err := api.auth.GenerateToken(kid, claims)

	if err != nil {
		return errs.New(errs.Internal, "failed to generate token")
	}

	resp := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}

	return web.Respond(ctx, w, resp, http.StatusOK)
}
