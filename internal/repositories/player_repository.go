package repositories

import (
	"scattergories-backend/config"
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/models"

	"gorm.io/gorm"
)

func GetPlayersByGameID(gameID uint) ([]*models.Player, error) {
	var players []*models.Player
	if err := config.DB.Where("game_id = ?", gameID).Find(&players).Error; err != nil {
		return nil, err
	}
	return players, nil
}

func GetPlayerByUserIDGameID(userID uint, gameID uint) (*models.Player, error) {
	var player models.Player
	if err := config.DB.Preload("User").Where("user_id = ? AND game_id = ?", userID, gameID).First(&player).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrPlayerNotFound
		}
		return nil, err
	}
	return &player, nil
}

func CreatePlayer(gamePlayer *models.Player) error {
	return config.DB.Create(gamePlayer).Error
}
