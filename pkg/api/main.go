package api

import (
	"net/http"
	"time"

	"go.deepl.dev/mealie-webhook-handler/pkg/appcontext"
)

func NewServer(appCtx appcontext.AppContext) http.Server {
	http.HandleFunc("/webhook/{identifier}", CreateHandleWebhook(appCtx))
	http.HandleFunc("/", HandleNotFound)

	srv := http.Server{
		Addr:         ":2025",
		WriteTimeout: 10 * time.Second,
	}
	return srv
}
