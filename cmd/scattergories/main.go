package main

import (
	"encoding/json"
	"log"
	"os"
	"scattergories-backend/config"
	"scattergories-backend/internal/api/routes"
	"scattergories-backend/internal/api/websocket"
	"scattergories-backend/internal/rabbitmq"
	"scattergories-backend/pkg/validators"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	validators.RegisterCustomValidators()

	appConfig, err := config.InitializeApp()
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	// Defer the closure of connections until main exits
	defer appConfig.DB.Close()
	defer appConfig.RedisClient.Close()
	defer appConfig.RabbitMQ.Close()

	go consumeLLMResponses(appConfig)

	go websocket.HubInstance.Run()

	router := gin.Default()
	routes.RegisterRoutes(
		router,
		appConfig.AuthService,
		appConfig.TokenService,
		appConfig.GameRoomService,
		appConfig.GameRoomJoinService,
		appConfig.UserService,
		appConfig.UserRegistrationService,
		appConfig.PermissionService,
		appConfig.MessageHandler,
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}

func consumeLLMResponses(appConfig *config.AppConfig) {
	msgs, err := appConfig.RabbitMQ.Consume("llm_response_queue")
	if err != nil {
		log.Fatalf("Failed to consume LLM response messages: %v", err)
	}

	for msg := range msgs {
		processLLMResponseMessage(msg.Body, appConfig)
	}
}

func processLLMResponseMessage(message []byte, appConfig *config.AppConfig) {
	var responseMsg rabbitmq.ResponseMessage
	err := json.Unmarshal(message, &responseMsg)
	if err != nil {
		log.Printf("Error unmarshaling LLM response message: %v", err)
		return
	}

	// Use the GameID to process the response and store the results in the database
	err = appConfig.AnswerValidationService.ProcessLLMResponse(responseMsg.GameID, responseMsg.Response)
	if err != nil {
		log.Printf("Error processing LLM response: %v", err)
		return
	}

	// Notify the client via WebSocket
	err = appConfig.MessageHandler.NotifyClientsLLMCompleted(uuid.MustParse(responseMsg.GameID))
	if err != nil {
		log.Printf("Error notifying clients: %v", err)
	}
}
