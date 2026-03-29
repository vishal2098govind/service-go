package checkapi

import (
	"github.com/vishal2098govind/service/app/api/auth"
	"github.com/vishal2098govind/service/foundations/logger"
	"github.com/vishal2098govind/service/foundations/web"
)

func Routes(app *web.App, log *logger.Logger, ath *auth.Auth) {

	app.HandleFunc("GET /liveness", liveness)
	app.HandleFunc("GET /readiness", readiness)
	app.HandleFunc("GET /testerror", testerror)
	app.HandleFunc("GET /testpanic", testpanic)
}
