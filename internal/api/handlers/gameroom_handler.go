package handlers

import (
	"net/http"
	"scattergories-backend/internal/api/handlers/requests"
	"scattergories-backend/internal/api/handlers/responses"
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func GetAllGameRooms(c *gin.Context) {
	rooms, err := services.GetAllGameRooms()
	if err != nil {
		HandleError(c, http.StatusInternalServerError, "Failed to retrieve game rooms")
		return
	}

	response := make([]*responses.GameRoomResponse, 0, len(rooms))
	for _, room := range rooms {
		response = append(response, responses.ToGameRoomResponse(room))
	}

	c.JSON(http.StatusOK, response)
}

func GetGameRoom(c *gin.Context) {
	roomID, err := GetIDParam(c, "room_id")
	if err != nil {
		HandleError(c, http.StatusBadRequest, "Invalid room ID")
		return
	}

	userID, _ := c.Get("userID")
	permitted, err := services.HasPermission(userID.(uint), services.GameRoomReadPermission, roomID)
	if err != nil || !permitted {
		HandleError(c, http.StatusUnauthorized, err.Error())
		return
	}

	room, err := services.GetGameRoomByID(roomID)
	if err != nil {
		if err == common.ErrGameRoomNotFound {
			HandleError(c, http.StatusNotFound, err.Error())
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
		if err == common.ErrUserNotFound {
			HandleError(c, http.StatusNotFound, err.Error())
		} else if err == common.ErrUserIsAlreadyHostOfAnotherRoom {
			HandleError(c, http.StatusConflict, err.Error())
		} else {
			HandleError(c, http.StatusInternalServerError, "Failed to create game room")
		}
		return
	}

	response := responses.ToGameRoomResponse(gameRoom)
	c.JSON(http.StatusCreated, response)
}

func DeleteGameRoom(c *gin.Context) {
	id, err := GetIDParam(c, "room_id")
	if err != nil {
		HandleError(c, http.StatusBadRequest, "Invalid room ID")
		return
	}

	userID, _ := c.Get("userID")
	permitted, err := services.HasPermission(userID.(uint), services.GameRoomWritePermission, id)
	if err != nil || !permitted {
		HandleError(c, http.StatusUnauthorized, err.Error())
		return
	}

	err = services.DeleteGameRoomByID(id)
	if err != nil {
		if err == common.ErrGameRoomNotFound {
			HandleError(c, http.StatusNotFound, err.Error())
		} else {
			HandleError(c, http.StatusInternalServerError, "Failed to delete game room")
		}
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Room deleted"})
}
