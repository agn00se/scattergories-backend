package routes

import (
	"scattergories-backend/internal/services"

	"scattergories-backend/internal/api/handlers"
	"scattergories-backend/internal/api/websocket"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.Engine,
	authService services.AuthService,
	tokenService services.TokenService,
	gameRoomService services.GameRoomService,
	gameRoomJoinService services.GameRoomJoinService,
	userService services.UserService,
	userRegistrationService services.UserRegistrationService,
	permissionService services.PermissionService,
	messageHandler websocket.MessageHandler,
) {
	authHandler := handlers.NewAuthHandler(authService, tokenService)
	userHandler := handlers.NewUserHandler(userService, userRegistrationService, tokenService, permissionService)
	gameRoomHandler := handlers.NewGameRoomHandler(gameRoomService, permissionService)
	gameRoomJoinHandler := handlers.NewGameRoomJoinHandler(gameRoomJoinService)

	RegisterUserRoutes(router, authHandler, userHandler, tokenService)
	RegisterGameRoomRoutes(router, gameRoomHandler, gameRoomJoinHandler, tokenService, gameRoomService, messageHandler)
}
