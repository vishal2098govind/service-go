package mux

import (
	"os"

	"github.com/vishal2098govind/service/apis/services/api/middlewares"
	"github.com/vishal2098govind/service/apis/services/sales/route/sys/checkapi"
	"github.com/vishal2098govind/service/app/api/auth"
	"github.com/vishal2098govind/service/foundations/logger"
	"github.com/vishal2098govind/service/foundations/web"
)

func WebAPI(log *logger.Logger, auth *auth.Auth, shutdown chan os.Signal) *web.App {
	mux := web.NewApp(shutdown, middlewares.Logger(log), middlewares.Errors(log), middlewares.Panics(log))

	checkapi.Routes(mux, log, auth)

	return mux
}
