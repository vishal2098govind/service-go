package checkapi

import (
	"github.com/vishal2098govind/service/foundations/web"
)

func Routes(app *web.App) {
	app.HandleFunc("GET /liveness", liveness)
	app.HandleFunc("GET /readiness", readiness)
}
