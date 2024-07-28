package routes

import (
	"scattergories-backend/internal/api/handlers"
	"scattergories-backend/internal/api/middlewares"
	"scattergories-backend/internal/api/websocket"
	"scattergories-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterGameRoomRoutes(
	router *gin.Engine,
	gameRoomHandler handlers.GameRoomHandler,
	gameRoomJoinHandler handlers.GameRoomJoinHandler,
	tokenService services.TokenService,
	gameRoomService services.GameRoomService,
	messageHandler websocket.MessageHandler,
) {
	gameRoomRoutes := router.Group("/game-rooms")
	gameRoomRoutes.Use(middlewares.JWTAuthMiddleware(tokenService))
	{
		gameRoomRoutes.GET("", gameRoomHandler.GetAllGameRooms)
		gameRoomRoutes.GET("/:room_id", gameRoomHandler.GetGameRoom)
		gameRoomRoutes.POST("", gameRoomHandler.CreateGameRoom)
		gameRoomRoutes.DELETE("/:room_id", gameRoomHandler.DeleteGameRoom) // might not need to be exposed as API endpoint

		gameRoomRoutes.PUT("/:room_id/join", gameRoomJoinHandler.JoinGameRoom)
		gameRoomRoutes.PUT("/:room_id/leave", gameRoomJoinHandler.LeaveGameRoom)
	}

	router.GET("/ws/:room_id", middlewares.JWTAuthMiddleware(tokenService), func(c *gin.Context) {
		websocket.HandleWebSocket(c, messageHandler)
	})
}
