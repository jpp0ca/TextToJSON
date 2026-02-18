package http

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// NewRouter registers all routes and returns an http.Handler.
func NewRouter(handler *StructureHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /structure", handler.StructureText)

	// Swagger UI
	mux.Handle("GET /swagger/", httpSwagger.WrapHandler)

	return mux
}
