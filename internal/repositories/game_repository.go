package repositories

import (
	"scattergories-backend/config"
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/domain"

	"gorm.io/gorm"
)

func GetGameByID(id uint) (*domain.Game, error) {
	var game domain.Game
	if err := config.DB.First(&game, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrGameNotFound
		}
		return nil, err
	}
	return &game, nil
}

func GetGameByRoomIDAndStatus(roomID uint, statuses ...string) (*domain.Game, error) {
	var game domain.Game
	err := config.DB.Where("game_room_id = ? AND status IN ?", roomID, statuses).First(&game).Error
	if err != nil {
		return nil, err
	}
	return &game, nil
}

func CreateGame(game *domain.Game) error {
	return config.DB.Create(game).Error

}

func UpdateGame(game *domain.Game) error {
	return config.DB.Save(game).Error
}
