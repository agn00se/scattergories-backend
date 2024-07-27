package routes

import (
	"scattergories-backend/internal/api/handlers"
	"scattergories-backend/internal/api/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine) {
	userRoutes := router.Group("/users")
	{
		userRoutes.POST("", handlers.CreateAccount) // Public route

		userRoutes.Use(middlewares.JWTAuthMiddleware())
		{
			userRoutes.GET("", handlers.GetAllUsers)
			userRoutes.GET("/:id", handlers.GetUser)
			// userRoutes.PUT("/:id", controllers.UpdateUser)
			userRoutes.DELETE("/:id", handlers.DeleteAccount)
		}
	}

	router.POST("/guests", handlers.CreateGuestAccount) // Public route
	router.POST("/login", handlers.Login)               // Public route
	router.POST("/logout", middlewares.JWTAuthMiddleware(), handlers.Logout)
	router.POST("/refresh-token", handlers.ExchangeToken) // Public route
}
