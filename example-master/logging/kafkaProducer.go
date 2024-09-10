package logging

import (
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"golang.org/x/example/model"
)

var (
	_kafkaServer string
	_kafkaTopic  string
)

// SetKafkaLoggingConfig initializes Kafka server and topic from configuration.
func SetKafkaLoggingConfig(c model.Configs) {
	_kafkaServer = c.KafkaServer
	_kafkaTopic = c.KafkaTopic
	log.Printf("Kafka server set to: %s", _kafkaServer)
	log.Printf("Kafka topic set to: %s", _kafkaTopic)
}

// sendUniqueIDCountToKafka sends the unique ID count to a Kafka topic.
func sendUniqueIDCountToKafka(count int64) error {
	// Create a new Kafka producer
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": _kafkaServer,
	})
	if err != nil {
		log.Printf("Failed to create Kafka producer with server %s: %v", _kafkaServer, err)
		return err
	}
	defer producer.Close()

	// Prepare the message to send
	message := fmt.Sprintf("Unique ID count: %d", count)
	log.Printf("Preparing to send message to Kafka topic %s: %s", _kafkaTopic, message)

	// Produce the message to the Kafka topic
	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &_kafkaTopic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, nil)
	if err != nil {
		log.Printf("Failed to produce message to Kafka topic %s: %v", _kafkaTopic, err)
		return err
	}

	// Wait for message delivery
	producer.Flush(15 * 1000)
	log.Printf("Message:- %v successfully sent to Kafka topic %s", message, _kafkaTopic)
	return nil
}
