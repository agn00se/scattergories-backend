package routes

import (
	"scattergories-backend/internal/client/controllers"
	"scattergories-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterGameRoomRoutes(router *gin.Engine) {
	gameRoomRoutes := router.Group("/game-rooms")
	gameRoomRoutes.Use(middleware.JWTAuthMiddleware())
	{
		gameRoomRoutes.GET("", controllers.GetAllGameRooms)
		gameRoomRoutes.GET("/:room_id", controllers.GetGameRoom)
		gameRoomRoutes.POST("", controllers.CreateGameRoom)
		gameRoomRoutes.DELETE("/:room_id", controllers.DeleteGameRoom)

		gameRoomRoutes.PUT("/:room_id/join", controllers.JoinGameRoom)
		gameRoomRoutes.PUT("/:room_id/leave", controllers.LeaveGameRoom)
	}
}
