package routes

import (
	"log"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"golang.org/x/example/deduplication"
	"golang.org/x/example/logging"
)

// Handler struct for your handlers
type Handler struct{}

// CorsMiddleware applies CORS settings
func (h *Handler) CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// AcceptGetHandler Handler
func (h *Handler) AcceptGetHandler(w http.ResponseWriter, r *http.Request) {
	// Record the start time of request processing
	startTime := time.Now()

	// Extract 'id' and 'endpoint' query parameters from the request URL
	idStr := r.URL.Query().Get("id")
	endpoint := r.URL.Query().Get("endpoint")

	// Check if 'id' parameter is missing
	if idStr == "" {
		http.Error(w, "Missing 'id' query parameter", http.StatusBadRequest)
		log.Printf("Error: Missing 'id' query parameter from request URL: %s", r.URL.String())
		return
	}

	// Convert 'id' from string to integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid 'id' query parameter", http.StatusBadRequest)
		log.Printf("Error: Invalid 'id' query parameter '%s' from request URL: %s", idStr, r.URL.String())
		return
	}

	// Store the ID and update the unique count if the ID is unique
	if deduplication.StoreUniqueID(id) {
		atomic.AddInt64(&logging.UniqueCount, 1)
		log.Printf("Stored unique ID: %d, updated unique count", id)
	} else {
		log.Printf("Duplicate ID: %d ignored", id)
	}

	// If an 'endpoint' parameter is provided, send asynchronous requests to it
	if endpoint != "" {
		log.Printf("Sending asynchronous requests to endpoint: %s", endpoint)
		go logging.SendGETRequestToEndpoint(endpoint)
		go logging.SendPOSTRequestToEndpoint(endpoint)
	} else {
		log.Printf("No 'endpoint' query parameter provided; skipping endpoint requests")
	}

	_, err = w.Write([]byte("ok"))
	if err != nil {
		log.Println("Failed to write the final Response with error :- ", err.Error())
		return
	}
	logProcessingTime(startTime, "Request processed successfully")
}

// logProcessingTime logs the elapsed time since startTime with a custom message
func logProcessingTime(startTime time.Time, message string) {
	elapsedTime := time.Since(startTime)
	log.Printf("%s, processing time: %v", message, elapsedTime)
}
