package mux

import (
	"os"

	"github.com/vishal2098govind/service/apis/services/api/mid"
	"github.com/vishal2098govind/service/apis/services/sales/route/sys/checkapi"
	"github.com/vishal2098govind/service/foundations/logger"
	"github.com/vishal2098govind/service/foundations/web"
)

func WebAPI(log *logger.Logger, shutdown chan os.Signal) *web.App {
	mux := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Panics(log))

	checkapi.Routes(mux)

	return mux
}
