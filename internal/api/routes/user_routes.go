package routes

import (
	"scattergories-backend/internal/api/handlers"
	"scattergories-backend/internal/api/middlewares"
	"scattergories-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine, authHandler handlers.AuthHandler, userHandler handlers.UserHandler, tokenService services.TokenService) {
	userRoutes := router.Group("/users")
	{
		userRoutes.POST("", userHandler.CreateAccount) // Public route

		userRoutes.Use(middlewares.JWTAuthMiddleware(tokenService))
		{
			userRoutes.GET("", userHandler.GetAllUsers)
			userRoutes.GET("/:id", userHandler.GetUser)
			// userRoutes.PUT("/:id", controllers.UpdateUser)
			userRoutes.DELETE("/:id", userHandler.DeleteAccount)
		}
	}

	router.POST("/guests", userHandler.CreateGuestAccount) // Public route
	router.POST("/login", authHandler.Login)               // Public route
	router.POST("/logout", middlewares.JWTAuthMiddleware(tokenService), authHandler.Logout)
	router.POST("/refresh-token", authHandler.ExchangeToken) // Public route
}
