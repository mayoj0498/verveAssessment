package main

import (
	"log"
)

func main() {
	// Start the GoVerve service and handle any errors
	err := startGoVerveService()
	if err != nil {
		log.Fatalf("Error starting GoVerve service: %v", err)
	}
}
