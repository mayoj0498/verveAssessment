package routes

import (
	"github.com/gorilla/mux"
	"net/http"
)

const (
	_verveAcceptRoute = "/api/verve/accept"
)

// SetupRoutes configures the routes for the server
func SetupRoutes(router *mux.Router) {
	// Create an instance of the handler struct
	h := &Handler{}

	router.HandleFunc(_verveAcceptRoute, h.AcceptGetHandler).Methods(http.MethodGet)

	// Optional: Set up middleware for CORS, logging, etc.
	router.Use(h.CorsMiddleware)
}
