package logging

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync/atomic"
)

// SendGETRequestToEndpoint sends an HTTP request with unique request count.
func SendGETRequestToEndpoint(endpoint string) {
	uniqueCount := int(atomic.LoadInt64(&UniqueCount))
	uniqueCountStr := strconv.Itoa(uniqueCount)

	// Create URL object
	baseURL, err := url.Parse(endpoint)
	if err != nil {
		log.Printf("Error parsing URL: %s", err)
		return
	}

	// Create query parameters
	params := url.Values{}
	params.Add("uniqueCount", uniqueCountStr)

	// Add the encoded query parameters to the URL
	baseURL.RawQuery = params.Encode()

	// Create a POST request with jsonData as the body
	resp, err := http.Get(baseURL.String())
	if err != nil {
		log.Printf("Failed to send GET request to endpoint: %v", err)
		return
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Printf("GET request - failed to close response body")
		}
	}(resp.Body)
	requestLog.Printf("GET request sent to endpoint: %v, returned status code: %d", endpoint, resp.StatusCode)
}

// SendPOSTRequestToEndpoint sends an HTTP request with unique request count.
func SendPOSTRequestToEndpoint(endpoint string) {
	count := atomic.LoadInt64(&UniqueCount)

	// Prepare JSON body for POST request
	data := map[string]int64{
		"uniqueCount": count,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal data: %v", err)
		return
	}

	// Create a POST request with jsonData as the body
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Failed to send POST request to endpoint: %v", err)
		return
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Printf("POST request - failed to close response body")
		}
	}(resp.Body)
	requestLog.Printf("POST request sent to endpoint: %v, returned status code: %d", endpoint, resp.StatusCode)
}
