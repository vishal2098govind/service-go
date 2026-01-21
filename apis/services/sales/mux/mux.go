package mux

import (
	"os"

	"github.com/vishal2098govind/service/apis/services/sales/route/sys/checkapi"
	"github.com/vishal2098govind/service/foundations/web"
)

func WebAPI(shutdown chan os.Signal) *web.App {
	mux := web.NewApp(shutdown)

	checkapi.Routes(mux)

	return mux
}
