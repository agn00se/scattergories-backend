package main

import (
	"log"
	"scattergories-backend/config"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize worker app configuration
	workerAppConfig, err := config.InitializeWorkerApp()
	if err != nil {
		log.Fatalf("Failed to initialize worker app: %v", err)
	}

	defer workerAppConfig.RabbitMQ.Close()

	// Consume messages from RabbitMQ and process them
	msgs, err := workerAppConfig.RabbitMQ.Consume("llm_validation_queue")
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	log.Println("Worker started. Waiting for messages...")

	// Process messages
	for msg := range msgs {
		log.Printf("Received message: %s", msg.Body)
		// 	err := services.HandleMessage(msg.Body)
		// if err != nil {
		// 		log.Printf("Error processing message: %v", err)
		// 	}
	}
}
