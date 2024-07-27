package main

import (
	"scattergories-backend/config"
	"scattergories-backend/internal/api/routes"
	"scattergories-backend/internal/api/websocket"
	"scattergories-backend/pkg/validators"

	"github.com/gin-gonic/gin"
)

func main() {

	validators.RegisterCustomValidators()

	config.ConnectDB()
	config.InitRedis()
	config.LoadPrompts()

	go websocket.HubInstance.Run()

	router := gin.Default()
	routes.RegisterRoutes(router)

	router.Run(":8080")
}
