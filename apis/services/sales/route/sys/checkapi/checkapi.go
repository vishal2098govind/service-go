package checkapi

import (
	"encoding/json"
	"net/http"
)

func liveness(w http.ResponseWriter, r *http.Request) {
	resp := struct {
		Status string
	}{
		Status: "OK",
	}

	json.NewEncoder(w).Encode(resp)
}

func readiness(w http.ResponseWriter, r *http.Request) {
	resp := struct {
		Status string
	}{
		Status: "OK",
	}

	json.NewEncoder(w).Encode(resp)
}
