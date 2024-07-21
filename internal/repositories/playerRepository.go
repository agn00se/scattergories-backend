package repositories

import (
	"scattergories-backend/config"
	"scattergories-backend/internal/models"
)

func CreatePlayer(gamePlayer *models.Player) error {
	return config.DB.Create(gamePlayer).Error
}
