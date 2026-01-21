package web

import (
	"context"
	"fmt"
	"net/http"
	"os"
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

	// mw := func(handler Handler) Handler {
	// 	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// 		// middle-ware specific
	// 		err := handler(ctx, w, r)
	// 		// middle-ware specific
	// 		return err
	// 	}
	// }

	// handler = mw(handler)

	h := func(w http.ResponseWriter, r *http.Request) {

		// CAN PUT SOME CODE HERE

		if err := handler(r.Context(), w, r); err != nil {
			// CAN HANDLER ERROR HERE
			fmt.Println(err)
		}

		// CAN PUT SOME CODE HERE

	}

	a.ServeMux.HandleFunc(pattern, h)
}
