package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"golang.org/x/example/cmd/webapp/routes"
	"golang.org/x/example/logging"
	"golang.org/x/example/model"
)

// startGoVerveService initializes and starts the GoVerve service
func startGoVerveService() error {
	// Load and set the configuration
	var configs model.Configs
	configs = setConfig()
	logging.SetKafkaLoggingConfig(configs)
	log.Println("Logging configuration set")

	// Set up HTTP routes
	routes.SetupRoutes()
	log.Println("Routes have been set up")

	// Start the HTTP server
	if err := startServer(configs); err != nil {
		return err
	}

	// Start the routine that logs unique requests every minute
	go logging.LogUniqueRequestsEveryMinute()
	log.Println("Started logging routine")
	return nil
}

// setConfig loads configuration from a YAML file and unmarshals it into a Configs struct
func setConfig() model.Configs {
	log.Println("Loading configuration from settings.yaml")

	// Read the YAML configuration file
	yamlFile, err := os.ReadFile("settings.yaml")
	if err != nil {
		log.Printf("Error reading YAML file: %v", err)
		log.Fatalf("Failed to read configuration file: %v", err)
	}

	// Unmarshal YAML data into Configs struct
	configs := model.Configs{}
	err = yaml.Unmarshal(yamlFile, &configs)
	if err != nil {
		log.Fatalf("Error unmarshalling YAML data: %v", err)
	}

	log.Println("Configuration successfully loaded")
	return configs
}

// startServer initializes and starts the HTTP server with the given configuration
func startServer(c model.Configs) error {
	// Configure the HTTP server
	server := &http.Server{
		Addr:         c.Host,
		ReadTimeout:  time.Duration(c.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(c.WriteTimeout) * time.Second,
	}

	// Log the server start
	log.Printf("Starting server on address: %s", c.Host)

	// Start the server and handle any errors
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Error while starting server: %v", err)
		log.Fatalf("Server failed: %v", err)
		return err
	}
	return nil
}
