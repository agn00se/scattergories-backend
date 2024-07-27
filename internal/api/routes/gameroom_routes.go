package routes

import (
	"scattergories-backend/internal/api/handlers"
	"scattergories-backend/internal/api/middlewares"
	"scattergories-backend/internal/api/websocket"

	"github.com/gin-gonic/gin"
)

func RegisterGameRoomRoutes(router *gin.Engine) {
	gameRoomRoutes := router.Group("/game-rooms")
	gameRoomRoutes.Use(middlewares.JWTAuthMiddleware())
	{
		gameRoomRoutes.GET("", handlers.GetAllGameRooms)
		gameRoomRoutes.GET("/:room_id", handlers.GetGameRoom)
		gameRoomRoutes.POST("", handlers.CreateGameRoom)
		gameRoomRoutes.DELETE("/:room_id", handlers.DeleteGameRoom) // might not need to be exposed as API endpoint

		gameRoomRoutes.PUT("/:room_id/join", handlers.JoinGameRoom)
		gameRoomRoutes.PUT("/:room_id/leave", handlers.LeaveGameRoom)
	}

	router.GET("/ws/:room_id", middlewares.JWTAuthMiddleware(), websocket.HandleWebSocket)
}
