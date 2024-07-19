package controllers

import (
	"net/http"
	"scattergories-backend/internal/client/controllers/requests"
	"scattergories-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func JoinGameRoom(c *gin.Context) {
	roomID, err := GetIDParam(c, "room_id")
	if err != nil {
		HandleError(c, http.StatusBadRequest, "Invalid room ID")
		return
	}

	var request requests.JoinLeaveRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		HandleError(c, http.StatusBadRequest, "Invalid request")
		return
	}

	err = services.JoinGameRoom(request.UserID, roomID)
	if err != nil {
		if err == services.ErrGameRoomNotFound {
			HandleError(c, http.StatusNotFound, "Game room not found")
		} else if err == services.ErrActiveGameExists {
			HandleError(c, http.StatusConflict, "Active game exists in the room")
		} else if err == services.ErrUserNotFound {
			HandleError(c, http.StatusNotFound, "User not found")
		} else {
			HandleError(c, http.StatusInternalServerError, "Failed to join game room")
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User joined game room"})
}

func LeaveGameRoom(c *gin.Context) {
	roomID, err := GetIDParam(c, "room_id")
	if err != nil {
		HandleError(c, http.StatusBadRequest, "Invalid room ID")
		return
	}

	var request requests.JoinLeaveRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		HandleError(c, http.StatusBadRequest, "Invalid request")
		return
	}

	err = services.LeaveGameRoom(request.UserID, roomID)
	if err != nil {
		if err == services.ErrGameRoomNotFound {
			HandleError(c, http.StatusNotFound, "Game room not found")
		} else if err == services.ErrUserNotFound {
			HandleError(c, http.StatusNotFound, "User not found")
		} else if err == services.ErrUserNotInSpecifiedRoom {
			HandleError(c, http.StatusNotFound, "User is not in the specified game room")
		} else {
			HandleError(c, http.StatusInternalServerError, "Failed to leave game room")
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User left game room"})
}
