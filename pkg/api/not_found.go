package api

import "net/http"

func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	SendError(w, http.StatusNotFound, "Not found")
}
