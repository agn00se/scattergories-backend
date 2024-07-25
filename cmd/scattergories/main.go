package main

import (
	"scattergories-backend/config"
	"scattergories-backend/internal/client/routes"
	"scattergories-backend/internal/client/ws"
	"scattergories-backend/pkg/validators"

	"github.com/gin-gonic/gin"
)

func main() {

	validators.RegisterCustomValidators()

	config.ConnectDB()
	config.InitRedis()
	config.LoadPrompts()

	go ws.HubInstance.Run()

	router := gin.Default()
	routes.RegisterRoutes(router)

	router.Run(":8080")
}
