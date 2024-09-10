package routes

import (
	"net/http"
)

const (
	_verveAcceptRoute = "/api/verve/accept"
)

// SetupRoutes configures the routes for the server
func SetupRoutes() {
	http.HandleFunc(_verveAcceptRoute, AcceptHandler)
}
