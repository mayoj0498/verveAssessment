package deduplication

import (
	"log"
	"sync"
)

var (
	uniqueIDs = make(map[int]bool)
	mu        sync.Mutex
)

// StoreUniqueID stores the ID if it is unique.
// Returns true if the ID was unique and false if it was a duplicate.
func StoreUniqueID(id int) bool {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := uniqueIDs[id]; exists {
		log.Printf("ID %d is a duplicate and was not stored.", id)
		return false
	}

	uniqueIDs[id] = true
	log.Printf("ID %d stored successfully.", id)
	return true
}

// ClearUniqueIDs resets the unique ID map to clear all stored IDs.
// This function should be called periodically to ensure the map is cleared.
func ClearUniqueIDs() {
	mu.Lock()
	defer mu.Unlock()

	// Log the clearing action for traceability
	log.Printf("Clearing all unique IDs. Previous count: %d", len(uniqueIDs))
	uniqueIDs = make(map[int]bool)
}
