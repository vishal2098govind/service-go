package checkapi

import (
	"github.com/vishal2098govind/service/apis/services/api/middlewares"
	"github.com/vishal2098govind/service/app/api/auth"
	"github.com/vishal2098govind/service/foundations/logger"
	"github.com/vishal2098govind/service/foundations/web"
)

func Routes(app *web.App, log *logger.Logger, ath *auth.Auth) {
	authenticate := middlewares.Authenticate(log, ath)
	ruleAdminOnly := auth.RuleAdminOnly
	authorizeAdminOnly := middlewares.Authorize(log, ath, ruleAdminOnly)

	app.HandleFunc("GET /liveness", liveness, authenticate, authorizeAdminOnly)
	app.HandleFunc("GET /readiness", readiness)
	app.HandleFunc("GET /testerror", testerror)
	app.HandleFunc("GET /testpanic", testpanic)
}
