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

func GetGameIDByGamePromptID(gamePromptID uint) (uint, error) {
	var gamePrompt models.GamePrompt
	if err := config.DB.Where("id = ?", gamePromptID).First(&gamePrompt).Error; err != nil {
		return 0, err
	}
	return gamePrompt.GameID, nil
}

func CreateGamePrompt(gamePrompt *models.GamePrompt) error {
	return config.DB.Create(gamePrompt).Error
}
