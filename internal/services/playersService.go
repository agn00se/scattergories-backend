package services

import (
	"scattergories-backend/config"
	"scattergories-backend/internal/models"
)

// Only doing one query search instead of checking if game_id, room_id player_id existed

func GetPlayerByID(roomID uint, gameID uint, playerID uint) (models.Player, error) {
	var player models.Player
	if err := config.DB.Where("id = ? AND game_id = ? AND game_room_id = ?", playerID, gameID, roomID).First(&player).Error; err != nil {
		return models.Player{}, err
	}
	return player, nil
}

func GetPlayersByGameID(roomID uint, gameID uint) ([]models.Player, error) {
	var players []models.Player
	if err := config.DB.Where("game_id = ? AND game_room_id = ?", gameID, roomID).Find(&players).Error; err != nil {
		return players, err
	}
	return players, nil
}

// GetAllPlayers
