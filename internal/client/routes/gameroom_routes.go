package routes

import (
	"scattergories-backend/internal/client/controllers"
	"scattergories-backend/internal/client/ws"
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
		gameRoomRoutes.DELETE("/:room_id", controllers.DeleteGameRoom) // might not need to be exposed as API endpoint

		gameRoomRoutes.PUT("/:room_id/join", controllers.JoinGameRoom)
		gameRoomRoutes.PUT("/:room_id/leave", controllers.LeaveGameRoom)
	}

	router.GET("/ws/:room_id", middleware.JWTAuthMiddleware(), ws.HandleWebSocket)
}
