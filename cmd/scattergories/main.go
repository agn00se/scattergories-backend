package main

import (
	"scattergories-backend/config"
	"scattergories-backend/internal/client/routes"
	"scattergories-backend/pkg/validators"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	validators.RegisterCustomValidators()

	config.ConnectDB()
	routes.RegisterRoutes(router)

	router.Run(":8080")
}
