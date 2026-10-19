package httpapi

import (
	"log/slog"
	"net/http"

	"catalog/internal/catalog"
)

// NewRouter returns the HTTP handler of the catalog service.
func NewRouter(store *catalog.Store, logger *slog.Logger) http.Handler {
	h := &handlers{store: store}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", h.index)
	mux.HandleFunc("GET /products/{id}", h.product)
	mux.HandleFunc("GET /healthz", h.healthz)

	return logRequests(logger, recoverPanics(logger, mux))
}
