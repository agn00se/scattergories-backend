package handlers

import (
	"net/http"
	"scattergories-backend/internal/api/handlers/requests"
	"scattergories-backend/internal/api/handlers/responses"
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type GameRoomHandler interface {
	GetAllGameRooms(c *gin.Context)
	GetGameRoom(c *gin.Context)
	CreateGameRoom(c *gin.Context)
	DeleteGameRoom(c *gin.Context)
}

type GameRoomHandlerImpl struct {
	gameRoomService   services.GameRoomService
	permissionService services.PermissionService
}

func NewGameRoomHandler(gameRoomService services.GameRoomService, permissionService services.PermissionService) GameRoomHandler {
	return &GameRoomHandlerImpl{gameRoomService: gameRoomService, permissionService: permissionService}
}

func (h *GameRoomHandlerImpl) GetAllGameRooms(c *gin.Context) {
	rooms, err := h.gameRoomService.GetAllGameRooms()
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

func (h *GameRoomHandlerImpl) GetGameRoom(c *gin.Context) {
	roomID, err := GetIDParam(c, "room_id")
	if err != nil {
		HandleError(c, http.StatusBadRequest, "Invalid room ID")
		return
	}

	userID, _ := c.Get("userID")
	permitted, err := h.permissionService.HasPermission(userID.(uint), services.GameRoomReadPermission, roomID)
	if err != nil || !permitted {
		HandleError(c, http.StatusUnauthorized, err.Error())
		return
	}

	room, err := h.gameRoomService.GetGameRoomByID(roomID)
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

func (h *GameRoomHandlerImpl) CreateGameRoom(c *gin.Context) {
	var request requests.GameRoomRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		HandleError(c, http.StatusBadRequest, err.Error())
		return
	}

	hostID := c.MustGet("userID").(uint)

	gameRoom, err := h.gameRoomService.CreateGameRoom(hostID, request.IsPrivate, request.Passcode)
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

func (h *GameRoomHandlerImpl) DeleteGameRoom(c *gin.Context) {
	id, err := GetIDParam(c, "room_id")
	if err != nil {
		HandleError(c, http.StatusBadRequest, "Invalid room ID")
		return
	}

	userID, _ := c.Get("userID")
	permitted, err := h.permissionService.HasPermission(userID.(uint), services.GameRoomWritePermission, id)
	if err != nil || !permitted {
		HandleError(c, http.StatusUnauthorized, err.Error())
		return
	}

	err = h.gameRoomService.DeleteGameRoomByID(id)
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
