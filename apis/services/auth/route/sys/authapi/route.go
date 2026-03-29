package authapi

import (
	"github.com/vishal2098govind/service/apis/services/api/middlewares"
	"github.com/vishal2098govind/service/app/api/auth"
	"github.com/vishal2098govind/service/foundations/logger"
	"github.com/vishal2098govind/service/foundations/web"
)

func Routes(mux *web.App, log *logger.Logger, auth *auth.Auth) {
	authApi := newAPI(auth, log)

	basic := middlewares.Basic(log, auth)

	mux.HandleFunc("POST /token/{kid}", authApi.generateToken, basic)
}
