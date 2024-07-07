package controllers

import (
	"net/http"
	"scattergories-backend/internal/client/controllers/requests"
	"scattergories-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func JoinGameRoom(c *gin.Context) {
	roomID, err := getIDParam(c, "room_id")
	if err != nil {
		handleError(c, http.StatusBadRequest, "Invalid room ID")
		return
	}

	var request requests.JoinLeaveRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid request")
		return
	}

	err = services.JoinGameRoom(request.UserID, roomID)
	if err != nil {
		if err == services.ErrGameRoomNotFound {
			handleError(c, http.StatusNotFound, "Game room not found")
		} else if err == services.ErrActiveGameExists {
			handleError(c, http.StatusConflict, "Active game exists in the room")
		} else if err == services.ErrUserNotFound {
			handleError(c, http.StatusNotFound, "User not found")
		} else {
			handleError(c, http.StatusInternalServerError, "Failed to join game room")
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User joined game room"})
}

func LeaveGameRoom(c *gin.Context) {
	roomID, err := getIDParam(c, "room_id")
	if err != nil {
		handleError(c, http.StatusBadRequest, "Invalid room ID")
		return
	}

	var request requests.JoinLeaveRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid request")
		return
	}

	err = services.LeaveGameRoom(request.UserID, roomID)
	if err != nil {
		if err == services.ErrGameRoomNotFound {
			handleError(c, http.StatusNotFound, "Game room not found")
		} else if err == services.ErrUserNotFound {
			handleError(c, http.StatusNotFound, "User not found")
		} else if err == services.ErrUserNotInSpecifiedRoom {
			handleError(c, http.StatusBadRequest, "User is not in the specified game room")
		} else {
			handleError(c, http.StatusInternalServerError, "Failed to leave game room")
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User left game room"})
}
