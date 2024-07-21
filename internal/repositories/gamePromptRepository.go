package repositories

import (
	"scattergories-backend/config"
	"scattergories-backend/internal/models"
)

func GetGamePromptsByGameID(gameID uint) ([]*models.GamePrompt, error) {
	var gamePrompts []*models.GamePrompt
	if err := config.DB.Where("game_id = ?", gameID).Preload("Prompt").Find(&gamePrompts).Error; err != nil {
		return nil, err
	}
	return gamePrompts, nil
}

func CreateGamePrompt(gamePrompt *models.GamePrompt) error {
	return config.DB.Create(gamePrompt).Error
}
