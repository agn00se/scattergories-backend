package handlers

import (
	"net/http"
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type GameRoomJoinHandler interface {
	JoinGameRoom(c *gin.Context)
	LeaveGameRoom(c *gin.Context)
}

type GameRoomJoinHandlerImpl struct {
	gameRoomJoinService services.GameRoomJoinService
}

func NewGameRoomJoinHandler(gameRoomJoinService services.GameRoomJoinService) GameRoomJoinHandler {
	return &GameRoomJoinHandlerImpl{gameRoomJoinService: gameRoomJoinService}
}

func (h *GameRoomJoinHandlerImpl) JoinGameRoom(c *gin.Context) {
	roomID, err := GetIDParam(c, "room_id")
	if err != nil {
		HandleError(c, http.StatusBadRequest, "Invalid room ID")
		return
	}

	userID := c.MustGet("userID").(uint)

	err = h.gameRoomJoinService.JoinGameRoom(userID, roomID)
	if err != nil {
		if err == common.ErrGameRoomNotFound || err == common.ErrUserNotFound {
			HandleError(c, http.StatusNotFound, err.Error())
		} else if err == common.ErrActiveGameExists {
			HandleError(c, http.StatusConflict, err.Error())
		} else {
			HandleError(c, http.StatusInternalServerError, "Failed to join game room")
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User joined game room"})
}

func (h *GameRoomJoinHandlerImpl) LeaveGameRoom(c *gin.Context) {
	roomID, err := GetIDParam(c, "room_id")
	if err != nil {
		HandleError(c, http.StatusBadRequest, "Invalid room ID")
		return
	}

	userID := c.MustGet("userID").(uint)

	err = h.gameRoomJoinService.LeaveGameRoom(userID, roomID)
	if err != nil {
		if err == common.ErrGameRoomNotFound || err == common.ErrUserNotFound || err == common.ErrUserNotInSpecifiedRoom {
			HandleError(c, http.StatusNotFound, err.Error())
		} else if err == common.ErrUserIsAlreadyHostOfAnotherRoom {
			HandleError(c, http.StatusConflict, err.Error())
		} else {
			HandleError(c, http.StatusInternalServerError, "Failed to leave game room")
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User left game room"})
}
