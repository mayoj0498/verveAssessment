package logging

import (
	"log"
	"os"
	"time"
)

const (
	TimeFormat = "2006-01-02_15-04-05"
)

// writeToLogFile writes the unique request count to a log file.
func writeToLogFile(count int64) error {
	// Generate a timestamped log file name
	currentTime := time.Now().Format(TimeFormat)
	logFileName := "logfile_" + currentTime + ".log"
	log.Printf("Opening log file: %s", logFileName)

	// Open the log file with append and write permissions
	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		// Log an error if the file cannot be opened
		log.Printf("Failed to open log file %s: %v", logFileName, err)
		return err
	}
	defer func() {
		// Close the file after logging
		if cerr := file.Close(); cerr != nil {
			log.Printf("Error closing log file %s: %v", logFileName, cerr)
		}
	}()

	// Set the log output to the file
	log.SetOutput(file)

	// Log the unique request count to the file
	requestLog.Printf("Unique request count : %d", count)
	log.Printf("Successfully logged unique request count %d to file %s", count, logFileName)

	return nil
}
