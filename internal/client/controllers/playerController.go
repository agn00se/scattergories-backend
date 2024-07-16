package controllers

import (
	"net/http"
	"scattergories-backend/internal/client/controllers/responses"
	"scattergories-backend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetPlayer(c *gin.Context) {
	roomID, err := GetIDParam(c, "room_id")
	if err != nil {
		HandleError(c, http.StatusBadRequest, "Invalid room ID")
		return
	}

	gameID, err := GetIDParam(c, "game_id")
	if err != nil {
		HandleError(c, http.StatusBadRequest, "Invalid game ID")
		return
	}

	playerID, err := GetIDParam(c, "player_id")
	if err != nil {
		HandleError(c, http.StatusBadRequest, "Invalid player ID")
		return
	}

	player, err := services.GetPlayerByID(roomID, gameID, playerID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			HandleError(c, http.StatusNotFound, "Player not found")
		} else {
			HandleError(c, http.StatusInternalServerError, "Failed to get player")
		}
		return
	}

	response := responses.ToPlayerResponse(player)
	c.JSON(http.StatusOK, response)
}

func GetPlayersByGameID(c *gin.Context) {
	roomID, err := GetIDParam(c, "room_id")
	if err != nil {
		HandleError(c, http.StatusBadRequest, "Invalid room ID")
		return
	}

	gameID, err := GetIDParam(c, "game_id")
	if err != nil {
		HandleError(c, http.StatusBadRequest, "Invalid game ID")
		return
	}

	players, err := services.GetPlayersByGameID(roomID, gameID)
	if err != nil {
		HandleError(c, http.StatusInternalServerError, "Failed to retrieve players")
		return
	}

	var response []responses.PlayerResponse
	for _, player := range players {
		response = append(response, responses.ToPlayerResponse(player))
	}

	c.JSON(http.StatusOK, response)
}
