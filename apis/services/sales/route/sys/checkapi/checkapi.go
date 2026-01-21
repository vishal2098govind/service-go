package checkapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/vishal2098govind/service/foundations/web"
)

func liveness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	resp := struct {
		Status string
	}{
		Status: "OK",
	}

	return web.Respond(ctx, w, resp, http.StatusOK)
}

func readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	resp := struct {
		Status string
	}{
		Status: "OK",
	}

	return web.Respond(ctx, w, resp, http.StatusOK)
}

func testerror(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return fmt.Errorf("error")
}
func testpanic(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	panic("panic")
}
