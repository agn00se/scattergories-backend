package repositories

import (
	"scattergories-backend/config"
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/models"

	"gorm.io/gorm"
)

func GetPlayerByID(id uint) (*models.Player, error) {
	var player models.Player
	if err := config.DB.Preload("User").First(&player, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrPlayerNotFound
		}
		return nil, err
	}
	return &player, nil
}

func GetPlayersByGameID(gameID uint) ([]*models.Player, error) {
	var players []*models.Player
	if err := config.DB.Where("game_id = ?", gameID).Find(&players).Error; err != nil {
		return nil, err
	}
	return players, nil
}

func CreatePlayer(gamePlayer *models.Player) error {
	return config.DB.Create(gamePlayer).Error
}
