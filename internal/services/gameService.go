package services

import (
	"scattergories-backend/config"
	"scattergories-backend/internal/models"
	"scattergories-backend/pkg/utils"

	"time"
)

func CreateGame(roomID uint) (models.Game, error) {
	// Check if the GameRoom has any ongoing or voting games
	var existingGames []models.Game
	if err := config.DB.Where("game_room_id = ? AND (status = ? OR status = ?)", roomID, models.GameStatusOngoing, models.GameStatusVoting).Find(&existingGames).Error; err != nil {
		return models.Game{}, err
	}

	if len(existingGames) > 0 {
		return models.Game{}, ErrActiveGameExists
	}

	// Create the new game with the status set to Ongoing
	game := models.Game{
		GameRoomID: roomID,
		Status:     models.GameStatusOngoing,
		StartTime:  time.Now(),
		Letter:     utils.GetRandomLetter(),
	}
	if err := config.DB.Create(&game).Error; err != nil {
		return game, err
	}

	// Find all users in the GameRoom and create Player entries for the new game
	var users []models.User
	if err := config.DB.Where("game_room_id = ?", roomID).Find(&users).Error; err != nil {
		return game, err
	}

	for _, user := range users {
		gamePlayer := models.Player{
			UserID: user.ID,
			GameID: game.ID,
			Score:  0,
		}
		if err := config.DB.Create(&gamePlayer).Error; err != nil {
			return game, err
		}
	}

	return game, nil
}

func GetGamesByRoomID(roomID uint) ([]models.Game, error) {
	var games []models.Game
	if err := config.DB.Where("game_room_id = ?", roomID).Find(&games).Error; err != nil {
		return games, err
	}
	return games, nil
}

func GetGameByID(roomID uint, gameID uint) (models.Game, error) {
	var game models.Game
	if err := config.DB.Where("id = ? AND game_room_id = ?", gameID, roomID).First(&game).Error; err != nil {
		return models.Game{}, err
	}
	return game, nil
}
