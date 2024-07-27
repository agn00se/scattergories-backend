package repositories

import (
	"scattergories-backend/config"
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/domain"

	"gorm.io/gorm"
)

func GetPlayersByGameID(gameID uint) ([]*domain.Player, error) {
	var players []*domain.Player
	if err := config.DB.Where("game_id = ?", gameID).Find(&players).Error; err != nil {
		return nil, err
	}
	return players, nil
}

func GetPlayerByUserIDGameID(userID uint, gameID uint) (*domain.Player, error) {
	var player domain.Player
	if err := config.DB.Preload("User").Where("user_id = ? AND game_id = ?", userID, gameID).First(&player).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrPlayerNotFound
		}
		return nil, err
	}
	return &player, nil
}

func CreatePlayer(gamePlayer *domain.Player) error {
	return config.DB.Create(gamePlayer).Error
}
