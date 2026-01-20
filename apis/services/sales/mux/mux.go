package mux

import (
	"encoding/json"
	"net/http"
)

func WebAPI() *http.ServeMux {
	mux := http.NewServeMux()

	h := func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			Status string
		}{
			Status: "OK",
		}

		json.NewEncoder(w).Encode(resp)
	}

	mux.HandleFunc("/test", h)

	return mux
}
