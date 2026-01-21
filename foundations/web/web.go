package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type App struct {
	*http.ServeMux
	shutdown chan os.Signal
	mids     []MidHandler
}

func NewApp(shutdown chan os.Signal, mids ...MidHandler) *App {
	return &App{
		ServeMux: http.NewServeMux(),
		shutdown: shutdown,
		mids:     mids,
	}
}

func (a *App) HandleFunc(pattern string, handler Handler, mids ...MidHandler) {
	handler = wrapMiddlewares(handler, mids...)
	handler = wrapMiddlewares(handler, a.mids...)

	h := func(w http.ResponseWriter, r *http.Request) {

		// CAN PUT SOME CODE HERE
		// can include trace here to the context as this is from where the request starts
		// this is the first point of contact of the request with the service
		ctx := setValues(r.Context(), &Value{
			TraceID: uuid.NewString(),
			Now:     time.Now().UTC(),
		})

		// calling handler makes the incoming request pass through all the middlewares if any, and then pass through the original handler
		if err := handler(ctx, w, r); err != nil {
			// CAN HANDLER ERROR HERE
			fmt.Println(err)
		}

		// CAN PUT SOME CODE HERE

	}

	a.ServeMux.HandleFunc(pattern, h)
}
