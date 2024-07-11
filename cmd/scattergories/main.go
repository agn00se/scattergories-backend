package main

import (
	"scattergories-backend/config"
	"scattergories-backend/internal/client/routes"
	"scattergories-backend/internal/client/ws"
	"scattergories-backend/pkg/validators"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	validators.RegisterCustomValidators()

	config.ConnectDB()
	config.LoadPrompts()
	routes.RegisterRoutes(router)

	router.GET("/ws", func(c *gin.Context) {
		ws.HandleWebSocket(c)
	})

	router.Run(":8080")
}
