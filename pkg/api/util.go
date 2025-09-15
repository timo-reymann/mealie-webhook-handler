package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func SendError(w http.ResponseWriter, status int, msg string) {
	slog.Error("HTTP Error", "status", status, "message", msg)
	content, err := json.Marshal(Error{Message: msg})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(content)
}
