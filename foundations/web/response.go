package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func Respond(ctx context.Context, w http.ResponseWriter, data any, statusCode int) error {

	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("web.Respond: marshal: %w", err)
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)

	if _, err := w.Write(jsonData); err != nil {
		return fmt.Errorf("web.Respond: write: %w", err)
	}

	return nil
}
