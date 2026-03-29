package authapi

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/vishal2098govind/service/app/api/auth"
	"github.com/vishal2098govind/service/app/api/errs"
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

	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    api.auth.Issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Audience:  jwt.ClaimStrings{"sales"},
			Subject:   "4fc800d4-2c3d-45fc-a0fa-9263644f6de7",
		},
		Roles: []string{"USER"},
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
