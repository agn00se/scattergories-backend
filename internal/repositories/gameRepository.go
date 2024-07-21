package repositories

import (
	"scattergories-backend/config"
	"scattergories-backend/internal/models"
)

func GetGameByRoomIDAndStatus(roomID uint, statuses ...string) (*models.Game, error) {
	var game models.Game
	err := config.DB.Where("game_room_id = ? AND status IN ?", roomID, statuses).First(&game).Error
	if err != nil {
		return nil, err
	}
	return &game, nil
}

func CreateGame(game *models.Game) error {
	return config.DB.Create(game).Error

}

func UpdateGame(game *models.Game) error {
	return config.DB.Save(game).Error
}
