package controllers

import (
	"net/http"
	"scattergories-backend/internal/client/controllers/responses"
	"scattergories-backend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetPlayer(c *gin.Context) {
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

	playerID, err := getIDParam(c, "player_id")
	if err != nil {
		handleError(c, http.StatusBadRequest, "Invalid player ID")
		return
	}

	player, err := services.GetPlayerByID(roomID, gameID, playerID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(c, http.StatusNotFound, "Player not found")
		} else {
			handleError(c, http.StatusInternalServerError, "Failed to get player")
		}
		return
	}

	response := responses.ToPlayerResponse(player)
	c.JSON(http.StatusOK, response)
}

func GetPlayersByGameID(c *gin.Context) {
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

	players, err := services.GetPlayersByGameID(roomID, gameID)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to retrieve players")
		return
	}

	var response []responses.PlayerResponse
	for _, player := range players {
		response = append(response, responses.ToPlayerResponse(player))
	}

	c.JSON(http.StatusOK, response)
}
