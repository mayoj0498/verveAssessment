package logging

import (
	"log"
	"sync/atomic"
	"time"

	"golang.org/x/example/deduplication"
)

var (
	UniqueCount int64
	requestLog  = log.Default()
)

// LogUniqueRequestsEveryMinute logs the count of unique requests every minute
// and resets the counter and unique ID map.
func LogUniqueRequestsEveryMinute() {
	// Create a ticker that triggers every minute
	ticker := time.NewTicker(1 * time.Minute)
	for {
		// Wait for the ticker to tick
		<-ticker.C

		// Load the current unique count
		count := atomic.LoadInt64(&UniqueCount)
		log.Printf("Logging unique request count: %d", count)

		// Log the unique count to the log file with error handling
		if err := writeToLogFile(count); err != nil {
			log.Printf("Error logging to file: %v", err)
		} else {
			log.Printf("Successfully logged %d unique requests to file", count)
		}

		// Send the unique count to Kafka with error handling
		if err := sendUniqueIDCountToKafka(count); err != nil {
			log.Printf("Error sending unique count to Kafka: %v", err)
		} else {
			log.Printf("Successfully sent %d unique requests count to Kafka", count)
		}

		// Clear the unique ID map
		deduplication.ClearUniqueIDs()
		log.Println("Cleared unique ID map")

		// Reset the unique request counter
		atomic.StoreInt64(&UniqueCount, 0)
		log.Println("Reset unique request counter to 0")
	}
}
