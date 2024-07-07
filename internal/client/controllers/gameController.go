package controllers

import (
	"net/http"
	"scattergories-backend/internal/client/controllers/responses"
	"scattergories-backend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetGamesByRoomID(c *gin.Context) {
	roomID, err := getIDParam(c, "room_id")
	if err != nil {
		handleError(c, http.StatusBadRequest, "Invalid room ID")
		return
	}

	games, err := services.GetGamesByRoomID(roomID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(c, http.StatusNotFound, "Room not found")
		} else {
			handleError(c, http.StatusInternalServerError, "Failed to retrieve games")
		}
		return
	}

	var response []responses.GameResponse
	for _, game := range games {
		response = append(response, responses.ToGameResponse(game))
	}

	c.JSON(http.StatusOK, response)
}

func GetGame(c *gin.Context) {
	roomID, err := getIDParam(c, "room_id")
	if err != nil {
		handleError(c, http.StatusBadRequest, "Invalid room ID")
		return
	}

	gameID, err := getIDParam(c, "game_id")
	if err != nil {
		handleError(c, http.StatusBadRequest, "Invalid game ID")
		return
	}

	game, err := services.GetGameByID(roomID, gameID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(c, http.StatusNotFound, "Game not found")
		} else {
			handleError(c, http.StatusInternalServerError, "Failed to get game")
		}
		return
	}

	response := responses.ToGameResponse(game)
	c.JSON(http.StatusOK, response)
}

func CreateGame(c *gin.Context) {
	roomID, err := getIDParam(c, "room_id")
	if err != nil {
		handleError(c, http.StatusBadRequest, "Invalid room ID")
		return
	}

	game, err := services.CreateGame(roomID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(c, http.StatusNotFound, "Room not found")
		} else if err == services.ErrActiveGameExists {
			handleError(c, http.StatusConflict, err.Error())
		} else {
			handleError(c, http.StatusInternalServerError, "Failed to create game")
		}
		return
	}

	response := responses.ToGameResponse(game)
	c.JSON(http.StatusOK, response)
}
