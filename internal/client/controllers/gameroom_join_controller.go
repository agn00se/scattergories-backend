package controllers

import (
	"net/http"
	"scattergories-backend/internal/client/controllers/requests"
	"scattergories-backend/internal/common"
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
		HandleError(c, http.StatusBadRequest, err.Error())
		return
	}

	err = services.JoinGameRoom(request.UserID, roomID)
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

func LeaveGameRoom(c *gin.Context) {
	roomID, err := GetIDParam(c, "room_id")
	if err != nil {
		HandleError(c, http.StatusBadRequest, "Invalid room ID")
		return
	}

	var request requests.JoinLeaveRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		HandleError(c, http.StatusBadRequest, err.Error())
		return
	}

	err = services.LeaveGameRoom(request.UserID, roomID)
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
