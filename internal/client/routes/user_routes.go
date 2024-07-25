package routes

import (
	"scattergories-backend/internal/client/controllers"
	"scattergories-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine) {
	userRoutes := router.Group("/users")
	{
		userRoutes.POST("", controllers.CreateAccount) // Public route

		userRoutes.Use(middleware.JWTAuthMiddleware())
		{
			userRoutes.GET("", controllers.GetAllUsers)
			userRoutes.GET("/:id", controllers.GetUser)
			// userRoutes.PUT("/:id", controllers.UpdateUser)
			userRoutes.DELETE("/:id", controllers.DeleteAccount)
		}
	}

	router.POST("/guests", controllers.CreateGuestAccount) // Public route
	router.POST("/login", controllers.Login)               // Public route
	router.POST("/logout", middleware.JWTAuthMiddleware(), controllers.Logout)
	router.POST("/refresh-token", controllers.ExchangeToken) // Public route
}
