# _# **# Thought Process**_

## 1. Objective Understanding :-

* The goal of this project is to build a high-performance REST service, track unique requests based on an integer ID, and report these counts periodically. 
    - The service has one GET endpoint - /api/verve/accept which is able to accept an integer id as a
      mandatory query parameter and an optional string HTTP endpoint query parameter. It should return
      String “ok” if there were no errors processing the request and “failed” in case of any errors.
    - Every minute, the application should write the count of unique requests your application received in
      that minute to a log file - please use a standard logger. Uniqueness of request is based on the id
      parameter provided.
    - When the endpoint is provided, the service should fire an HTTP GET request to the provided endpoint
      with count of unique requests in the current minute as a query parameter. Also log the HTTP status
      code of the response.
* The implementation should also include the following extensions:
    - Perform an HTTP POST request as extension when an optional endpoint query parameter is provided.
    - Instead of logging the count of unique IDs to a log file, send it to a distributed streaming service (Uses Kafka).

## 2. Design Considerations :- 

#### To meet the requirement, the following design choices were made:

* Concurrency: Go's goroutines are used to process requests concurrently, enabling high throughput.
* Atomic Operations: The sync/atomic package is used to safely increment the unique request counter, ensuring that multiple goroutines can update the count without race conditions.
* Mutex for Deduplication: A mutex is used to protect access to the unique ID map (uniqueIDs) ensuring thread safety when storing IDs.

#### **Endpoint Flexibility**

The service has an endpoint `/api/verve/accept` that accepts two query parameters:
* id (int): Mandatory query parameter representing the unique ID.
* endpoint (string): Optional query parameter representing a URL to which an HTTP GET and POST request is fired with the unique request count.


## 3. Implementation Approach

#### REST Endpoint

The `/api/verve/accept` endpoint is implemented as a GET endpoint. The request-handling logic ensures:
* If an ID is received, the service checks if the ID is unique using the StoreUniqueID function.
* If an endpoint is provided, the service makes either a GET and POST request as Extension 1 to the provided URL. The HTTP request includes the current unique request count as a parameter (GET) or as JSON data (POST).

#### Periodic Unique Request Count Logging

Every minute, a goroutine logs the count of unique requests received in the last minute:
* Unique ID Count Logging: The unique count is written to a log file, and the unique IDs are cleared from memory after logging.
* Kafka Integration: Instead of just logging the count to a file, the count is also sent to a Kafka topic using the confluent-kafka-go library. This fulfills Extension 3.

#### Configuration

The service loads configuration details like Kafka settings, and server timeouts from a YAML file (`settings.yaml`). This ensures easy configurability for different environments.

## 4. Extensions

####    Extension 1: POST Request Instead of GET
   - The service supports firing both GET and POST request to the provided endpoint
   - The POST request sends the unique request count in the request body as JSON format.

####    Extension 3: Kafka Integration

   For scalability, instead of writing the unique ID count to a file, the application sends the count to a Kafka topic (unique-id-counts). This ensures the count is available for distributed systems or downstream services that can consume Kafka messages.

## 5. Key Functions and Modules

####    Main Application (main.go): Responsible for initializing configurations, setting up routes, and starting the server.

####    Deduplication Module (deduplication.go): Handles the deduplication logic, storing and clearing unique IDs.

####    Logging Module (logging package): Handles periodic logging of the unique request count to the log file and sending the count to Kafka.

####    HTTP Client (logging package): Sends HTTP GET and POST requests with the unique request count and log the status code of the request.

## 6. Error Handling and Logging

   Comprehensive error handling is in place across all operations, especially for external interactions such as writing to files, communicating with Kafka, and making HTTP requests.This ensures that any failures are logged for easier debugging and troubleshooting.

## 7. Conclusion

   This implementation balances high performance, scalability, and flexibility. By leveraging Go's concurrency features and integrating with Kafka, the service is designed to handle the expected high request load, ensure deduplication across multiple instances, and provide configurable, extensible logging mechanisms.