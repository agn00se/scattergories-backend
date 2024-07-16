package services

import (
	"scattergories-backend/config"
	"scattergories-backend/internal/models"
	"scattergories-backend/pkg/utils"
)

var (
	defaultTimeLimit       = 100
	defaultNumberOfPrompts = 10
)

func CreateDefaultGameRoomConfig(gameRoomID uint) error {
	gameRoomConfig := models.GameRoomConfig{
		GameRoomID:      gameRoomID,
		TimeLimit:       defaultTimeLimit,
		NumberOfPrompts: defaultNumberOfPrompts,
		Letter:          utils.GetRandomLetter(),
	}
	if err := config.DB.Create(&gameRoomConfig).Error; err != nil {
		return err
	}
	return nil
}
