package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/IBM/sarama"
)

type Message struct {
	EventType string    `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
	Priority  string    `json:"priority"`
	Source    string    `json:"source"`
	Location  string    `json:"location"`
}

var (
	brokers  = []string{"kafka:9092"} // Kafka broker addresses
	topic1    = "hello-world-topic-one"        // Kafka topic to produce to
	topic2    = "hello-world-topic-two"        // Kafka topic to produce to
	producer sarama.SyncProducer          // Kafka producer instance
)

func main() {
	// Initialize Kafka producer
	initKafkaProducer()

	// Create a Gin router
	router := gin.Default()

	// Define the POST endpoint
	router.POST("/message", handleMessage)

	// Start the HTTP server on port 8080
	fmt.Println("Server listening on port 8080...")
	router.Run(":8080")

	// Graceful shutdown handling
	defer func() {
		if err := producer.Close(); err != nil {
			log.Printf("Error closing Kafka producer: %v\n", err)
		}
	}()
}

func initKafkaProducer() {
	// Configure Kafka producer
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to respond
	// config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms
	config.Producer.Return.Successes = true                  // Enable success notifications

	// Initialize the producer
	var err error
	producer, err = sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}

	// Handle successful startup
	fmt.Println("Kafka producer initialized")
}

func handleMessage(c *gin.Context) {
	var msg Message
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Produce message to Kafka topic
	sendMessageToKafka(msg)

	// Send a response back
	c.JSON(http.StatusOK, gin.H{"message": "Message sent to Kafka successfully"})
}

func sendMessageToKafka(message Message) {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal message to JSON: %v\n", err)
		return
	}
	if message.EventType == "event1" {
		partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
			Topic: topic1,
			Value: sarama.StringEncoder(jsonMessage),
		})
		fmt.Printf("Message sent to partition %d at offset %d\n", partition, offset)
		if err != nil {
			log.Printf("Failed to send message to Kafka: %v\n", err)
			return
		}
	} else {
		partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
			Topic: topic2,
			Value: sarama.StringEncoder(jsonMessage),
		})
		fmt.Printf("Message sent to partition %d at offset %d\n", partition, offset)

		if err != nil {
			log.Printf("Failed to send message to Kafka: %v\n", err)
			return
		}
	}
}
