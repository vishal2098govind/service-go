package web

type MidHandler func(Handler) Handler

func wrapMiddlewares(handler Handler, mids ...MidHandler) Handler {

	for i := len(mids) - 1; i >= 0; i-- {
		mw := mids[i]
		if mw != nil {
			handler = mw(handler)
		}
	}

	return handler
}
