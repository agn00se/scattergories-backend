package routes

import (
	"scattergories-backend/internal/client/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterGameRoomRoutes(router *gin.Engine) {
	gameRoomRoutes := router.Group("/game-rooms")
	{
		gameRoomRoutes.GET("", controllers.GetAllGameRooms)
		gameRoomRoutes.GET("/:room_id", controllers.GetGameRoom)
		gameRoomRoutes.POST("", controllers.CreateGameRoom)
		gameRoomRoutes.DELETE("/:room_id", controllers.DeleteGameRoom)
		gameRoomRoutes.PUT("/:room_id/update-host", controllers.UpdateHost)

		gameRoomRoutes.POST("/:room_id/join", controllers.JoinGameRoom)
		gameRoomRoutes.POST("/:room_id/leave", controllers.LeaveGameRoom)
	}
}
