package handler

import (
	"log/slog"
	"net/http"
	"time"
)

// NewRouter registers all HTTP routes and wraps them in our standard middleware chain.
func NewRouter(logger *slog.Logger, requestTimeout time.Duration, health *HealthHandler , workshops *WorkshopHandler) http.Handler {
	mux := http.NewServeMux()

	// Operational probes
	mux.HandleFunc("GET /healthz", health.Healthz)
	mux.HandleFunc("GET /readyz", health.Readyz)
	// Workshop endpoints
	mux.HandleFunc("GET /v1/workshops", workshops.List)
	mux.HandleFunc("GET /v1/workshops/{id}", workshops.GetByID)
	mux.HandleFunc("POST /v1/workshops", workshops.Create)	
	// Build middleware chain from inside out:
	// Request -> RequestLogger -> Recover -> Timeout -> mux
	var h http.Handler = mux
	h = Timeout(requestTimeout, h)
	h = Recover(logger, h)
	h = RequestLogger(logger, h)

	return h
}