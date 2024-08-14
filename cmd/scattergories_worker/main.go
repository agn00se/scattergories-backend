package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"scattergories-backend/config"
	"scattergories-backend/internal/rabbitmq"

	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)

var openaiClient *openai.Client

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

	// Initialize OpenAI client
	openaiClient = openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	// Consume messages from RabbitMQ
	msgs, err := workerAppConfig.RabbitMQ.Consume("llm_request_queue")
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	log.Println("Worker started. Waiting for messages...")

	// Process messages
	for msg := range msgs {
		processMessage(workerAppConfig.RabbitMQ, msg.Body)
	}
}

func processMessage(rabbitMQ *rabbitmq.RabbitMQ, messageBody []byte) {
	var reqMsg rabbitmq.RequestMessage
	err := json.Unmarshal(messageBody, &reqMsg)
	if err != nil {
		log.Printf("Error unmarshaling request message: %v", err)
		return
	}

	response, err := sendLLMRequest(reqMsg.Prompt)
	if err != nil {
		log.Printf("Error sending LLM request: %v", err)
		return
	}

	responseMsg := rabbitmq.ResponseMessage{
		GameID:   reqMsg.GameID,
		Response: response,
	}
	responseBody, err := json.Marshal(responseMsg)
	if err != nil {
		log.Printf("Error marshaling response message: %v", err)
		return
	}

	err = rabbitMQ.Publish("llm_response_queue", responseBody)
	if err != nil {
		log.Printf("Error publishing LLM response message: %v", err)
	}
}

func sendLLMRequest(prompt string) (string, error) {
	response, err := openaiClient.CreateChatCompletion(
		context.TODO(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo0125,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens: 1000,
		},
	)
	if err != nil {
		return "", err
	}

	return response.Choices[0].Message.Content, nil
}
