package repositories

import (
	"scattergories-backend/config"
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/domain"

	"gorm.io/gorm"
)

func GetGameRoomConfigByRoomID(roomID uint) (*domain.GameRoomConfig, error) {
	var gameRoomConfig domain.GameRoomConfig
	if err := config.DB.First(&gameRoomConfig, "game_room_id = ?", roomID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrGameRoomConfigNotFound
		}
		return nil, err
	}
	return &gameRoomConfig, nil
}

func CreateGameRoomConfig(gameRoomConfig *domain.GameRoomConfig) error {
	return config.DB.Create(gameRoomConfig).Error
}

func UpdateGameRoomConfig(gameRoomConfig *domain.GameRoomConfig) error {
	return config.DB.Save(&gameRoomConfig).Error
}
