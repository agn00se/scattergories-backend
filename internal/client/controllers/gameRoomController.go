package controllers

import (
	"net/http"
	"scattergories-backend/internal/client/controllers/requests"
	"scattergories-backend/internal/client/controllers/responses"
	"scattergories-backend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllGameRooms(c *gin.Context) {
	rooms, err := services.GetAllGameRooms()
	if err != nil {
		HandleError(c, http.StatusInternalServerError, "Failed to retrieve game rooms")
	}

	var response []responses.GameRoomResponse
	for _, room := range rooms {
		response = append(response, responses.ToGameRoomResponse(room))
	}

	c.JSON(http.StatusOK, response)
}

func GetGameRoom(c *gin.Context) {
	id, err := GetIDParam(c, "room_id")
	if err != nil {
		HandleError(c, http.StatusBadRequest, "Invalid room ID")
		return
	}

	room, err := services.GetGameRoomByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			HandleError(c, http.StatusNotFound, "Room not found")
		} else {
			HandleError(c, http.StatusInternalServerError, "Failed to get game room")
		}
		return
	}

	response := responses.ToGameRoomResponse(room)
	c.JSON(http.StatusOK, response)
}

func CreateGameRoom(c *gin.Context) {
	var request requests.GameRoomRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		HandleError(c, http.StatusBadRequest, err.Error())
		return
	}

	gameRoom, err := services.CreateGameRoom(request.HostID, request.IsPrivate, request.Passcode)
	if err != nil {
		if err == services.ErrHostNotFound {
			HandleError(c, http.StatusNotFound, err.Error())
		} else if err == services.ErrUserIsAlreadyHostOfAnotherRoom {
			HandleError(c, http.StatusConflict, err.Error())
		} else {
			HandleError(c, http.StatusInternalServerError, "Failed to create game room")
		}
		return
	}

	response := responses.ToGameRoomResponse(gameRoom)
	c.JSON(http.StatusOK, response)
}

func DeleteGameRoom(c *gin.Context) {
	id, err := GetIDParam(c, "room_id")
	if err != nil {
		HandleError(c, http.StatusBadRequest, "Invalid room ID")
		return
	}

	err = services.DeleteGameRoomByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			HandleError(c, http.StatusNotFound, "Room not found")
		} else {
			HandleError(c, http.StatusInternalServerError, "Failed to delete game room")
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Room deleted"})
}

func UpdateHost(c *gin.Context) {
	id, err := GetIDParam(c, "room_id")
	if err != nil {
		HandleError(c, http.StatusBadRequest, "Invalid room ID")
		return
	}

	var request requests.UpdateHostRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		HandleError(c, http.StatusBadRequest, err.Error())
		return
	}

	gameRoom, err := services.UpdateHost(id, request.NewHostID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			HandleError(c, http.StatusNotFound, "Room not found")
		} else {
			HandleError(c, http.StatusInternalServerError, "Failed to update host")
		}
		return
	}

	response := responses.ToGameRoomResponse(gameRoom)
	c.JSON(http.StatusOK, response)
}
