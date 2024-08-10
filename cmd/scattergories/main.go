package main

import (
	"log"
	"os"
	"scattergories-backend/config"
	"scattergories-backend/internal/api/routes"
	"scattergories-backend/internal/api/websocket"
	"scattergories-backend/pkg/validators"

	"github.com/gin-gonic/gin"
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
