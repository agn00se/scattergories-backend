package repositories

import (
	"scattergories-backend/config"
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/models"

	"gorm.io/gorm"
)

func GetGameRoomConfigByRoomID(roomID uint) (*models.GameRoomConfig, error) {
	var gameRoomConfig models.GameRoomConfig
	if err := config.DB.First(&gameRoomConfig, "game_room_id = ?", roomID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrGameRoomConfigNotFound
		}
		return nil, err
	}
	return &gameRoomConfig, nil
}

func CreateGameRoomConfig(gameRoomConfig *models.GameRoomConfig) error {
	return config.DB.Create(gameRoomConfig).Error
}

func UpdateGameRoomConfig(gameRoomConfig *models.GameRoomConfig) error {
	return config.DB.Save(&gameRoomConfig).Error
}
