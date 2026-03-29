package mux

import (
	"net/http"
	"os"

	"github.com/vishal2098govind/service/apis/services/api/middlewares"
	"github.com/vishal2098govind/service/apis/services/auth/route/sys/authapi"
	"github.com/vishal2098govind/service/app/api/auth"
	"github.com/vishal2098govind/service/foundations/logger"
	"github.com/vishal2098govind/service/foundations/web"
)

func WebAPI(shutdown chan os.Signal, log *logger.Logger, auth *auth.Auth) http.Handler {
	mux := web.NewApp(shutdown, middlewares.Logger(log), middlewares.Errors(log), middlewares.Panics(log))

	authapi.Routes(mux, log, auth)

	return mux
}
